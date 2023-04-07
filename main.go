package main

import (
	"github.com/jeffadavidson/development-bot/logic/examinedata"
	"github.com/jeffadavidson/development-bot/utilities/config"
	"github.com/jeffadavidson/development-bot/utilities/exit"
)

func main() {
	config.ManualInit()

	err := examinedata.ExamineDevelopmentPermits()
	if err != nil {
		exit.ExitError(err)
	}

	exit.ExitSuccess()
}
