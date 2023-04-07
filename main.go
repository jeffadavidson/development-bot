package main

import (
	"fmt"

	"github.com/jeffadavidson/development-bot/utilities/config"
	"github.com/jeffadavidson/development-bot/utilities/exit"
)

//"github.com/jeffadavidson/development-bot/interactions/calgaryopendata"

func main() {
	config.ManualInit()

	//calgaryopendata.GetDevelopmentPermits()

	fmt.Println(config.Config.Neighborhood.Name)
	exit.ExitSuccess()
}
