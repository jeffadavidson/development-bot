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

	//Process all development activity into combined RSS feed
	err = examinedata.ProcessAllDevelopmentActivity()
	if err != nil {
		exit.ExitError(err)
	}
	fmt.Println("All Development Activity Processed Successfully")

	exit.ExitSuccess()
}

func ManualInits() error {
	configErr := config.ManualInit()
	if configErr != nil {
		return fmt.Errorf("Failed to start due to configuration error: %s", configErr.Error())
	}

	examineErr := examinedata.ManualInit()
	if examineErr != nil {
		return fmt.Errorf("Failed to start due to initialization error: %s", examineErr.Error())
	}

	return nil
}
