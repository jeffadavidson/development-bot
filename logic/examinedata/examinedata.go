package examinedata

import (
	"fmt"

	"github.com/jeffadavidson/development-bot/interactions/calgaryopendata"
	"github.com/jeffadavidson/development-bot/objects/developmentpermit"
	"github.com/jeffadavidson/development-bot/objects/fileaction"
	"github.com/jeffadavidson/development-bot/utilities/fileio"
	"github.com/jeffadavidson/development-bot/utilities/toolbox"
)

// ExamineDevelopmentPermits - Examines Development permits, determines actions to take
func ExamineDevelopmentPermits() error {
	fetchedDevelopmentPermits, storedDevelopmentPermits, err := getDevelopmentPermits()
	if err != nil {
		return err
	}

	fileActions := getDevelopmentPermitActions(fetchedDevelopmentPermits, storedDevelopmentPermits)

	for _, val := range fileActions {
		fmt.Println("------------------------------------------------")
		fmt.Println(val.PermitNum + ":")
		fmt.Println("\tAction: " + val.Action)
		fmt.Println("\tMessage: " + val.Message)
	}

	return nil
}

func isDevelopmentPermitClosed(fetchedDP developmentpermit.DevelopmentPermit, storedDP developmentpermit.DevelopmentPermit) (bool, string) {
	toClose := false
	closeMessage := ""

	// Check for close
	close_statuses := [3]string{"Released", "Cancelled", "Cancelled - Pending Refund"}
	if toolbox.SliceContains([]string(close_statuses[:]), fetchedDP.StatusCurrent) {
		toClose = true
		closeMessage = fmt.Sprintf("Closing file as it is in status '%s'", fetchedDP.StatusCurrent)
	}

	return toClose, closeMessage
}

func getDevelopmentPermitUpdates(fetchedDP developmentpermit.DevelopmentPermit, storedDP developmentpermit.DevelopmentPermit) (bool, string) {
	hasUpdate := false
	updateMessage := ""

	// check status
	if fetchedDP.StatusCurrent != storedDP.StatusCurrent {
		hasUpdate = true
		updateMessage += fmt.Sprintf("Status updated from '%s' to '%s'\n", storedDP.StatusCurrent, fetchedDP.StatusCurrent)
	}

	// check decision
	if fetchedDP.Decision != storedDP.Decision {
		hasUpdate = true
		updateMessage += fmt.Sprintf("Decision updated to '%s'\n", fetchedDP.Decision)
	}

	// check comment by date
	if fetchedDP.MustCommenceDate != storedDP.MustCommenceDate {
		hasUpdate = true
		updateMessage += fmt.Sprintf("Must Commence By Date updated to '%s'\n", fetchedDP.MustCommenceDate)
	}

	return hasUpdate, updateMessage
}

func getDevelopmentPermits() ([]developmentpermit.DevelopmentPermit, []developmentpermit.DevelopmentPermit, error) {
	// Load existing development permits
	storedDevelopmentPermitsBytes, loadErr := fileio.GetFileContents("./data/development-permits.json")
	if loadErr != nil {
		return nil, nil, loadErr
	}
	storedDevelopmentPermits, parseErr := calgaryopendata.ParseDevelopmentPermits(storedDevelopmentPermitsBytes)
	if parseErr != nil {
		return nil, nil, parseErr
	}

	//Get development Permits from calgary open data
	fetchedDevelopmentPermits, _ := calgaryopendata.GetDevelopmentPermits()

	return storedDevelopmentPermits, fetchedDevelopmentPermits, nil
}

func getDevelopmentPermitActions(fetchedDevelopmentPermits []developmentpermit.DevelopmentPermit, storedDevelopmentPermits []developmentpermit.DevelopmentPermit) []fileaction.FileAction {
	var fileActions []fileaction.FileAction
	for _, fetchedDP := range fetchedDevelopmentPermits {
		storedDpPtr := developmentpermit.FindDevelopmentPermit(storedDevelopmentPermits, fetchedDP.PermitNum)
		if storedDpPtr == nil {
			fileActions = append(fileActions, fileaction.FileAction{PermitNum: fetchedDP.PermitNum, Action: "CREATE", Message: fetchedDP.CreateInformationMessage()})
		} else {
			storedDP := *storedDpPtr

			// Skip if discussion closed
			if storedDP.GithubDiscussionClosed {
				fileActions = append(fileActions, fileaction.FileAction{PermitNum: fetchedDP.PermitNum, Action: "SKIP"})
				continue
			}

			hasUpdate, updateMessage := getDevelopmentPermitUpdates(fetchedDP, storedDP)
			toClose, closeMessage := isDevelopmentPermitClosed(fetchedDP, storedDP)

			message := fmt.Sprintf("%s:\n", storedDP.PermitNum)

			if hasUpdate && !toClose {
				message += updateMessage
				fileActions = append(fileActions, fileaction.FileAction{PermitNum: fetchedDP.PermitNum, Action: "UPDATE", Message: message})
			}
			if hasUpdate && toClose {
				message += updateMessage
				message += "\n"
				message += closeMessage

				fileActions = append(fileActions, fileaction.FileAction{PermitNum: fetchedDP.PermitNum, Action: "CLOSE", Message: message})
			}
			if !hasUpdate && toClose {
				message += closeMessage
				fileActions = append(fileActions, fileaction.FileAction{PermitNum: fetchedDP.PermitNum, Action: "CLOSE", Message: closeMessage})
			}
		}
	}

	return fileActions
}
