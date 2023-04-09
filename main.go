package main

import (
	"fmt"

	"github.com/jeffadavidson/development-bot/logic/examinedata"
	"github.com/jeffadavidson/development-bot/utilities/config"
	"github.com/jeffadavidson/development-bot/utilities/exit"
)

func main() {
	err := ManualInits()
	if err != nil {
		exit.ExitError(err)
	}

	//Trigger process for development permits
	dpErr := examinedata.DevelopmentPermits()
	if err != nil {
		exit.ExitError(dpErr)
	}

	exit.ExitSuccess()
}

func ManualInits() error {
	configErr := config.ManualInit()
	if configErr != nil {
		return fmt.Errorf("Failed to start due to configuration error: %s", configErr.Error())
	}
	examineErr := examinedata.ManualInit()
	if examineErr != nil {
		return fmt.Errorf("Failed to start due to date initialization error: %s", examineErr.Error())
	}

	return nil
}
