package examinedata

import (
	"fmt"

	"github.com/jeffadavidson/development-bot/objects/developmentpermit"
)

// ExamineDevelopmentPermits - Examines Development permits, determines actions to take
func ExamineDevelopmentPermits() error {
	fetchedDevelopmentPermits, storedDevelopmentPermits, err := developmentpermit.GetDevelopmentPermits()
	if err != nil {
		return err
	}

	fileActions := developmentpermit.GetDevelopmentPermitActions(fetchedDevelopmentPermits, storedDevelopmentPermits)

	for _, val := range fileActions {
		fmt.Println("------------------------------------------------")
		fmt.Println(val.PermitNum + ":")
		fmt.Println("\tAction: " + val.Action)
		fmt.Println("\tMessage: " + val.Message)
	}

	return nil
}
