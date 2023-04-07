package examinedata

import (
	"fmt"

	"github.com/jeffadavidson/development-bot/interactions/calgaryopendata"
	"github.com/jeffadavidson/development-bot/interactions/githubdiscussions"
	"github.com/jeffadavidson/development-bot/utilities/fileio"

	"golang.org/x/exp/slices"
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

	//var discussionsActions []githubdiscussions.DiscussionAction
	for _, fetchedDP := range fetchedDevelopmentPermits {
		var storedDP calgaryopendata.DevelopmentPermit

		//fmt.Printf("Index: %d, PermitNum %s, Status: %s\n", index, fetchedDP.PermitNum, fetchedDP.StatusCurrent)

		storedDPIndex := slices.IndexFunc(storedDevelopmentPermits, func(c calgaryopendata.DevelopmentPermit) bool { return c.PermitNum == fetchedDP.PermitNum })
		if storedDPIndex == -1 {
			//TODO:
			//fmt.Printf("Fresh Permit '%s' not found\n", fetchedDP.PermitNum)
		} else {
			storedDP = storedDevelopmentPermits[storedDPIndex]
			//fmt.Printf("Found Indexc: %d - %s\n", storedDPIndex, storedDP.PermitNum)

			if storedDP.GithubDiscussionClosed {
				//fmt.Printf("Skipping '%s' as discussion is closed\n", storedDP.PermitNum)
			} else {
				fmt.Printf("Checking '%s' for updates\n", storedDP.PermitNum)

				// compare status for update
				// compare decision for update
				// compare must comment by date for update
			}
		}

	}

	//itterate through all fresh development permits
	//find stored permit
	//if no find
	// create
	// if find
	// if archived, skip

	return nil
}
