package main

import (
	"fmt"

	"github.com/jeffadavidson/development-bot/interactions/calgaryopendata"
	"github.com/jeffadavidson/development-bot/utilities/config"
	"github.com/jeffadavidson/development-bot/utilities/exit"
)

func main() {
	config.ManualInit()

	developmentPermits, _ := calgaryopendata.GetDevelopmentPermits()
	fmt.Println(len(developmentPermits))

	exit.ExitSuccess()
}
