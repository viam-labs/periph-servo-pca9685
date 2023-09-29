// package main is a module that implements a servo model supported by a PCA9685 and periph.io
package main

import (
	"context"

	"github.com/edaniels/golog"
	"github.com/viam-labs/periph-servo-pca9685/pca9685"
	"go.viam.com/rdk/components/servo"
	"go.viam.com/rdk/module"
	"go.viam.com/utils"
)

func main() {
	utils.ContextualMain(mainWithArgs, golog.NewDevelopmentLogger("periph-servo-pca9685"))
}

func mainWithArgs(ctx context.Context, args []string, logger golog.Logger) error {
	servoModule, err := module.NewModuleFromArgs(ctx, logger)
	if err != nil {
		return err
	}

	servoModule.AddModelFromRegistry(ctx, servo.API, pca9685.Model)

	err = servoModule.Start(ctx)
	defer servoModule.Close(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}
