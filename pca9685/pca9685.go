//go:build linux

// Package pca9685servo implements a servo model supported by a PCA9685 and periph.io
package pca9685

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/edaniels/golog"
	"github.com/pkg/errors"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/devices/v3/pca9685"
	"periph.io/x/host/v3"

	"go.viam.com/rdk/components/servo"
	"go.viam.com/rdk/resource"
)

// Model represents a linux wifi strength sensor model.
var Model = resource.NewModel("viam-labs", "servo", "pca9685")

const defaultPwmFreq int = 50

type pca9685Servo struct {
	resource.Named
	resource.TriviallyCloseable
	resource.TriviallyReconfigurable
	logger golog.Logger

	servo                   *pca9685.Servo
	position                uint32
	cancelCtx               context.Context
	cancelFunc              func()
	activeBackgroundWorkers sync.WaitGroup
	mu                      sync.RWMutex
	moving                  bool
}

type Config struct {
	I2cBus           string `json:"i2c_bus"`
	I2cAddr          string `json:"i2c_addr"`
	Channel          int    `json:"channel"`
	Frequency        int    `json:"frequency_hz"`
	MinAngle         int    `json:"min_angle_deg"`
	MaxAngle         int    `json:"max_angle_deg"`
	StartingPosition uint32 `json:"starting_position_deg"`
	MinWidth         int    `json:"min_width_us"`
	MaxWidth         int    `json:"max_width_us"`
}

func (cfg *Config) Validate(path string) ([]string, error) {
	return []string{}, nil
}

func init() {
	resource.RegisterComponent(
		servo.API,
		Model,
		resource.Registration[servo.Servo, *Config]{
			Constructor: func(
				ctx context.Context,
				deps resource.Dependencies,
				conf resource.Config,
				logger golog.Logger,
			) (servo.Servo, error) {
				return newServo(ctx, deps, conf, logger)
			},
		})
}

func newServo(
	ctx context.Context,
	_ resource.Dependencies,
	conf resource.Config,
	logger golog.Logger,
) (servo.Servo, error) {
	logger.Info("Starting viam-labs:servo:pca9865 instance")

	cancelCtx, cancelFunc := context.WithCancel(context.Background())

	s := pca9685Servo{
		Named:      conf.ResourceName().AsNamed(),
		logger:     logger,
		position:   0,
		cancelCtx:  cancelCtx,
		cancelFunc: cancelFunc,
		mu:         sync.RWMutex{},
		moving:     false,
	}

	if err := s.Reconfigure(ctx, nil, conf); err != nil {
		return nil, err
	}
	return &s, nil
}

func (servo *pca9685Servo) Reconfigure(
	ctx context.Context,
	_ resource.Dependencies,
	conf resource.Config,
) error {
	servo.mu.Lock()
	defer servo.mu.Unlock()

	newConf, err := resource.NativeConfig[*Config](conf)
	if err != nil {
		return err
	}

	_, err = host.Init()
	if err != nil {
		return err
	}

	i2cBus := "0"
	if newConf.I2cBus != "" {
		i2cBus = newConf.I2cBus
	}

	bus, err := i2creg.Open(i2cBus)
	if err != nil {
		return err
	}

	i2cAddr := pca9685.I2CAddr
	if newConf.I2cAddr != "" {
		parsedAddr, err := strconv.ParseUint(newConf.I2cAddr, 0, 16)
		if err != nil {
			return err
		}
		i2cAddr = uint16(parsedAddr)
	}

	pca, err := pca9685.NewI2C(bus, i2cAddr)
	if err != nil {
		return err
	}

	pwmFreq := defaultPwmFreq
	if newConf.Frequency != 0 {
		pwmFreq = newConf.Frequency
	}

	if err := pca.SetPwmFreq(physic.Frequency(pwmFreq) * physic.Hertz); err != nil {
		return err
	}
	if err := pca.SetAllPwm(0, 0); err != nil {
		return err
	}

	// PWM cycle has 4096 time slots. Frequency hz * 4096 is time slots per second.
	// Divide that into 1 million to get microseconds per time slot
	usPerSlot := 1000000 / float64(pwmFreq*4096)

	minWidth := newConf.MinWidth
	if minWidth == 0 {
		minWidth = 500 // reasonable default based on common hobby servos
	}
	minPwm := int(float64(minWidth) / usPerSlot)

	maxWidth := newConf.MaxWidth
	if maxWidth == 0 {
		maxWidth = 2500 // reasonable default based on common hobby servos
	}
	maxPwm := int(float64(maxWidth) / usPerSlot)

	minAngle := 0
	if newConf.MinAngle != 0 {
		minAngle = newConf.MinAngle
	}
	maxAngle := 180
	if newConf.MaxAngle != 0 {
		maxAngle = newConf.MaxAngle
	}

	servo.logger.Info("Initializing servo on channel ", newConf.Channel, ", minPwm ", minPwm, ", maxPwm ", maxPwm, ", minAngle ", minAngle, ", maxAngle ", maxAngle, ", startingPosition ", newConf.StartingPosition)
	servos := pca9685.NewServoGroup(pca, gpio.Duty(minPwm), gpio.Duty(maxPwm), physic.Angle(minAngle), physic.Angle(maxAngle))

	servo.servo = servos.GetServo(newConf.Channel)

	servo.Move(ctx, newConf.StartingPosition, nil)

	return nil
}

// DoCommand always returns unimplemented but can be implemented by the embedder.
func (servo *pca9685Servo) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, resource.ErrDoUnimplemented
}

func (servo *pca9685Servo) IsMoving(ctx context.Context) (bool, error) {
	return servo.moving, nil
}

func (servo *pca9685Servo) Move(ctx context.Context, ang uint32, extra map[string]interface{}) error {
	servo.mu.Lock()
	defer servo.mu.Unlock()
	servo.moving = true
	if err := servo.servo.SetAngle(physic.Angle(ang)); err != nil {
		return errors.Wrap(err, "couldn't set angle")
	}
	servo.moving = false
	servo.position = ang
	return nil
}

func (servo *pca9685Servo) Position(ctx context.Context, extra map[string]interface{}) (uint32, error) {
	return uint32(servo.position), nil
}

func (servo *pca9685Servo) Stop(ctx context.Context, extra map[string]interface{}) error {
	servo.mu.Lock()
	defer servo.mu.Unlock()
	if err := servo.servo.SetPwm(0); err != nil {
		return errors.Wrap(err, "couldn't stop servo")
	}
	servo.moving = false
	return nil
}
