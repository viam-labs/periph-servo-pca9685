# `periph-servo-pca9685` modular component

This module implements the [Viam servo API](https://docs.viam.com/operate/reference/components/servo/#api) in a `periph-servo-pca9685` model.

With this model, you can control servos connected to PCA9685 channels using the [periph.io](https://periph.io/) library.

## Setup

Navigate to the **CONFIGURE** tab of your machine's page.

Click the **+** button, select **Component or service**, then select the `servo / periph-servo-pca9685` model provided by the module.

Click **Add module**, enter a name for your servo, and click **Create**.

## Configure your `periph-servo-pca9685` servo

On the new component panel, copy and paste the following attribute template into your servo's **Attributes** box:

```json
{
  "i2c_bus": <string>,
  "i2c_addr": <string>,
  "channel": <int>,
  "frequency_hz": <int>,
  "min_angle_deg": <int>,
  "max_angle_deg": <int>,
  "starting_position_deg": <int>,
  "min_width_us": <int>,
  "max_width_us": <int>
}
```

### Attributes

The following attributes are available for `periph-servo-pca9685` servos:

| Name | Type   | Inclusion | Description |
| ---- | ------ | --------- | ----------- |
| `i2c_bus` | string | Optional | The name or number of the I2C bus to which the PCA9685 is connected. Default: `"0"`. |
| `i2c_addr` | string | Optional | The I2C address of the PCA9685. Can be formatted as hex (0x) or base 10. See [this guide](https://learn.adafruit.com/scanning-i2c-addresses/raspberry-pi) for how to detect i2c devices. i2cdetect displays the hex formatted value. Default: ` "0x40"`. |
| `channel` | int | Optional | The channel (0-15) to which the servo is connected. Default: `0`. |
| `frequency_hz` | int | Optional | The frequency in Hz of the servo. See the servo datasheet. Default: `50`. |
| `min_angle_deg` | int | Optional | The minimum angle in degrees to which the servo will be allowed to move. Default: `0`. |
| `max_angle_deg` | int | Optional | The maximum angle in degrees to which the servo will be allowed to move. Default: `180`. |
| `starting_position_deg` | int | Optional | When the servo is initiated, it will move to this position (in degrees). Default: `0`. |
| `min_width_us` | int | Optional | The minimum duty cycle width in microseconds. See the servo datasheet. Default: `500`. |
| `max_width_us` | int | Optional | The maximum duty cycle width in microseconds. See the servo datasheet. Default: `2500`. |

### Example Configuration

```json
{
  "i2c_bus": "0",
  "channel": 15
}
```

## Development

To release a new version of this module, this repo uses the [Viam build-action](https://github.com/viamrobotics/build-action) to build the module in Viam's cloud infrastructure and deploy the new version based on a [tagged release](https://docs.github.com/en/repositories/releasing-projects-on-github/about-releases).

To kick off a deployment:

1. [Tag the release commit with the new module version](https://git-scm.com/book/en/v2/Git-Basics-Tagging) and push it to the repo
1. [Create a release based on that tag](https://docs.github.com/en/repositories/releasing-projects-on-github/managing-releases-in-a-repository#creating-a-release)

Within a couple of minutes, the new module version will be published to the Viam registry.

If there is an issue with the action and a manual release is required:

1. Authenticate the Viam CLI:

   ```console
   viam auth login
   ```

1. Start a remote build for the new module version

   ```console
   viam module build start --version <version>
   ```
