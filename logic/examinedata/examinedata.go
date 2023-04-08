package examinedata

import (
	"fmt"

	"github.com/jeffadavidson/development-bot/interactions/calgaryopendata"
	"github.com/jeffadavidson/development-bot/interactions/githubdiscussions"
	"github.com/jeffadavidson/development-bot/objects/developmentpermit"
	"github.com/jeffadavidson/development-bot/utilities/fileio"
	"github.com/jeffadavidson/development-bot/utilities/toolbox"
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
	fmt.Println(len(fetchedDevelopmentPermits))

	var discussionsActions []githubdiscussions.DiscussionAction
	for _, fetchedDP := range fetchedDevelopmentPermits {
		storedDP := developmentpermit.FindDevelopmentPermit(storedDevelopmentPermits, fetchedDP.PermitNum)
		if storedDP == nil {
			discussionsActions = append(discussionsActions, githubdiscussions.DiscussionAction{Action: "CREATE", Message: fetchedDP.CreateInformationMessage()})
		} else {
			//TODO: Make function to find development permit in array by permit num

			if storedDP.GithubDiscussionClosed {
				//fmt.Printf("Skipping '%s' as discussion is closed\n", storedDP.PermitNum)

			} else {
				fmt.Printf("Checking '%s' for updates\n", storedDP.PermitNum)
				hasUpdate := false
				updateMessage := ""

				// Check Status
				if storedDP.StatusCurrent != fetchedDP.StatusCurrent {
					hasUpdate = true
					updateMessage += "STATUS UPDATE "
				}

				// Check Decision
				if storedDP.Decision != fetchedDP.Decision {
					hasUpdate = true
					updateMessage += "DECITION UPDATE "
				}

				// Check Must Comment By Date
				if storedDP.MustCommenceDate != fetchedDP.MustCommenceDate {
					hasUpdate = true
					updateMessage += "MUST COMMENT BY UPDATE "
				}

				if hasUpdate {
					discussionsActions = append(discussionsActions, githubdiscussions.DiscussionAction{Action: "UPDATE", Message: updateMessage})
				}

				// Check for close
				close_discussion_statuses := [3]string{"Released", "Cancelled", "Cancelled - Pending Refund"}
				if toolbox.SliceContains([]string(close_discussion_statuses[:]), fetchedDP.StatusCurrent) {
					discussionsActions = append(discussionsActions, githubdiscussions.DiscussionAction{Action: "CLOSE", Message: "Close"})
				}

			}
		}

	}

	fmt.Println(discussionsActions)
	return nil
}
