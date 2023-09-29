//go:build linux

// Package pca9685servo implements a servo model supported by a PCA9685 and periph.io
package periph_pca9685

import (
	"context"
	"sync"

	"github.com/edaniels/golog"
	"github.com/pkg/errors"
	"go.viam.com/rdk/operation"

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

const defaultPwmFreq float32 = 50

type pca9685Servo struct {
	resource.Named
	resource.TriviallyCloseable
	resource.TriviallyReconfigurable
	logger golog.Logger
	opMgr  operation.SingleOperationManager

	servo                   *pca9685.Servo
	position                uint32
	cancelCtx               context.Context
	cancelFunc              func()
	activeBackgroundWorkers sync.WaitGroup
	mu                      sync.RWMutex
}

type Config struct {
	I2cBus           string `json:"i2c_bus"`
	Pin              int    `json:"pin"`
	Frequency        int    `json:"frequency_hz"`
	MinAngle         int    `json:"min_angle_deg"`
	MaxAngle         int    `json:"max_angle_deg"`
	StartingPosition int    `json:"starting_position_deg"`
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
		resource.Registration[servo.Servo, *Config]{Constructor: newServo})
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

	pca, err := pca9685.NewI2C(bus, pca9685.I2CAddr)
	if err != nil {
		return err
	}

	pwmFreq := 50
	if newConf.Frequency != 0 {
		pwmFreq = newConf.Frequency
	}

	if err := pca.SetPwmFreq(physic.Frequency(pwmFreq) * physic.Hertz); err != nil {
		return err
	}
	if err := pca.SetAllPwm(0, 0); err != nil {
		return err
	}

	minPwm := 50
	if newConf.MinWidth != 0 {
		minPwm = newConf.MinWidth
	}
	maxPwm := 650
	if newConf.MinWidth != 0 {
		maxPwm = newConf.MaxWidth
	}
	minAngle := 0
	if newConf.MinAngle != 0 {
		minAngle = newConf.MinAngle
	}
	maxAngle := 180
	if newConf.MaxAngle != 0 {
		maxAngle = newConf.MaxAngle
	}

	servos := pca9685.NewServoGroup(pca, gpio.Duty(minPwm), gpio.Duty(maxPwm), physic.Angle(minAngle), physic.Angle(maxAngle))

	// This is an example of using with an Me Arm robot arm
	servo.servo = servos.GetServo(0)

	return nil
}

// DoCommand always returns unimplemented but can be implemented by the embedder.
func (servo *pca9685Servo) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, resource.ErrDoUnimplemented
}

func (servo *pca9685Servo) IsMoving(ctx context.Context) (bool, error) {
	return servo.opMgr.OpRunning(), nil
}

func (servo *pca9685Servo) Move(ctx context.Context, ang uint32, extra map[string]interface{}) error {
	ctx, done := servo.opMgr.New(ctx)
	defer done()
	if err := servo.servo.SetAngle(physic.Angle(ang)); err != nil {
		return errors.Wrap(err, "couldn't set angle")
	}
	servo.position = ang
	return nil
}

func (servo *pca9685Servo) Position(ctx context.Context, extra map[string]interface{}) (uint32, error) {
	return uint32(servo.position), nil
}

func (servo *pca9685Servo) Stop(ctx context.Context, extra map[string]interface{}) error {
	ctx, done := servo.opMgr.New(ctx)
	defer done()
	// Turning the pin all the way off (i.e., setting the duty cycle to 0%) will cut power to the
	// motor. If you wanted to send it to position 0, you should set it to `minUs` instead.
	if err := servo.servo.SetPwm(0); err != nil {
		return errors.Wrap(err, "couldn't stop servo")
	}
	return nil
}
