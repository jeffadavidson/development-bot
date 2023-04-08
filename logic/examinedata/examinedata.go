package examinedata

import (
	"fmt"

	"github.com/jeffadavidson/development-bot/interactions/calgaryopendata"
	"github.com/jeffadavidson/development-bot/objects/developmentpermit"
	"github.com/jeffadavidson/development-bot/objects/fileaction"
	"github.com/jeffadavidson/development-bot/utilities/fileio"
)

// ExamineDevelopmentPermits - Produces a list of development permits where the stale and fresh state do not match
func ExamineDevelopmentPermits() error {
	// Load existing development permits
	storedDevelopmentPermitsBytes, loadErr := fileio.GetFileContents("./data/development-permits.json")
	if loadErr != nil {
		return loadErr
	}
	storedDevelopmentPermits, parseErr := calgaryopendata.ParseDevelopmentPermits(storedDevelopmentPermitsBytes)
	if parseErr != nil {
		return parseErr
	}

	//Get development Permits from calgary open data
	fetchedDevelopmentPermits, _ := calgaryopendata.GetDevelopmentPermits()

	var discussionsActions []fileaction.FileAction
	for _, fetchedDP := range fetchedDevelopmentPermits {
		storedDP := developmentpermit.FindDevelopmentPermit(storedDevelopmentPermits, fetchedDP.PermitNum)
		if storedDP == nil {
			discussionsActions = append(discussionsActions, fileaction.FileAction{PermitNum: fetchedDP.PermitNum, Action: "CREATE", Message: fetchedDP.CreateInformationMessage()})
		} else {

			// Skip if discussion closed
			if storedDP.GithubDiscussionClosed {
				discussionsActions = append(discussionsActions, fileaction.FileAction{PermitNum: fetchedDP.PermitNum, Action: "SKIP"})
				continue
			}

			//hasUpdate, updateMessage := GetDevelopmentPermitForUpdate()

			// fmt.Printf("Checking '%s' for updates\n", storedDP.PermitNum)
			// hasUpdate := false
			// updateMessage := ""

			// // Check Status
			// if storedDP.StatusCurrent != fetchedDP.StatusCurrent {
			// 	hasUpdate = true
			// 	updateMessage += "STATUS UPDATE "
			// }

			// // Check Decision
			// if storedDP.Decision != fetchedDP.Decision {
			// 	hasUpdate = true
			// 	updateMessage += "DECITION UPDATE "
			// }

			// // Check Must Comment By Date
			// if storedDP.MustCommenceDate != fetchedDP.MustCommenceDate {
			// 	hasUpdate = true
			// 	updateMessage += "MUST COMMENT BY UPDATE "
			// }

			// if hasUpdate {
			// 	discussionsActions = append(discussionsActions, fileaction.FileAction{Action: "UPDATE", Message: updateMessage})
			// }

			// // Check for close
			// close_discussion_statuses := [3]string{"Released", "Cancelled", "Cancelled - Pending Refund"}
			// if toolbox.SliceContains([]string(close_discussion_statuses[:]), fetchedDP.StatusCurrent) {
			// 	discussionsActions = append(discussionsActions, fileaction.FileAction{Action: "CLOSE", Message: "Close"})
			// }

		}

	}

	fmt.Println(discussionsActions)
	return nil
}
