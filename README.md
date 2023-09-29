# periph-servo-pca9685

A Viam modular servo implementation for servos connected to a PCA9685 breakout board

# periph-servo-pca9685

*periph-servo-pca9685* is a Viam modular component that uses the [periph.io](https://periph.io/) library to control servos connected to pca9685 channels.

## API

The periph-servo-pca9685 resource provides the following methods from Viam's built-in [rdk:component:servo API](https://docs.viam.com/components/servo/#api):

### Move(angle uint32)

### Position()

### Stop()

## Viam Component Configuration

The following attributes may be configured as facial-detector config attributes.
For example: the following configuration would use the `ssd` framework:

``` json
{
  	"i2c_bus" : "0",
	"channel" : 15,
	Frequency        int    `json:"frequency_hz"`
	MinAngle         int    `json:"min_angle_deg"`
	MaxAngle         int    `json:"max_angle_deg"`
	StartingPosition uint32 `json:"starting_position_deg"`
	MinWidth         int    `json:"min_width_us"`
	MaxWidth         int    `json:"max_width_us"`
}
```

### i2c_bus

*string (default: "0")*

The name or number of the I2C bus to which the PCA9685 is connected.

### channel

*int (default: 0)*

The channel (0-15) to which the servo to control is connected.

### frequency_hz

*int (default: 50)*

The frequency in hz of the servo to control.
See the datasheet for the servo to control.

### min_angle_deg

*int (default 0)*

The minimum angle in degrees to which the servo will be allowed to move.

### max_angle_deg

*int (default 180)*

The maximum angle to in degrees which the servo will be allowed to move.

### starting_position_deg

*int (default 0)*

When the servo is initiated, it will move to this position (in degrees).

### min_width_us

*int (default 500)*

The minimum duty cycle width in microseconds.
See the datasheet for the servo to control.
### max_width_us

*int (default 2500)*

The maximum duty cycle width in microseconds.
See the datasheet for the servo to control.
