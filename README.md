# periph-servo-pca9685

*periph-servo-pca9685* is a Viam modular component that uses the [periph.io](https://periph.io/) library to control servos connected to pca9685 channels.

## API

The periph-servo-pca9685 resource provides the following methods from Viam's built-in [rdk:component:servo API](https://docs.viam.com/components/servo/#api):

### Move(angle uint32)

### Position()

### Stop()

## Viam Component Configuration

The following attributes may be configured as periph-servo-pca9685 config attributes.
For example: the following configuration set up a servo on I2C bus 0, PCA9685 channel 15:

``` json
{
  "i2c_bus" : "0",
  "channel" : 15
}
```

### i2c_bus

*string (default: "0")*

The name or number of the I2C bus to which the PCA9685 is connected.

### i2c_addr

*string (default: "0x40")*

The number of the I2C address to which the PCA9685 is connected. This can be formatted as hex (prefixed by "0x") or base 10 (unprefixed) values.

If you're not sure which address to use, see [this guide] for how to detect i2c devices. `i2cdetect` displays the hex formatted value.

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

## Development

To release a new version of this module, this repo uses the [Viam build-action](https://github.com/viamrobotics/build-action) to build the module in Viam's cloud infrastructure and deploy the new version based on a [tagged release](https://docs.github.com/en/repositories/releasing-projects-on-github/about-releases).

To kick off a deployment:

1. [Tag the release commit with the new module version](https://git-scm.com/book/en/v2/Git-Basics-Tagging) and push it to the repo
1. [Create a release based on that tag](https://docs.github.com/en/repositories/releasing-projects-on-github/managing-releases-in-a-repository#creating-a-release)

Within a couple of minutes, the new module version should be published to the Viam registry.

If there is an issue with the action and a manual release is required:

1. Authenticate the Viam CLI:
   ```console
   viam auth login
   ```
1. Start a remote build for the new module version
   ```console
   viam module build start --version <version>
   ```
