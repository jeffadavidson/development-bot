package examinedata

import (
	//"fmt"

	"fmt"

	"github.com/jeffadavidson/development-bot/interactions/githubdiscussions"
	"github.com/jeffadavidson/development-bot/objects/developmentpermit"
	"github.com/jeffadavidson/development-bot/utilities/config"
)

// ExamineDevelopmentPermits - Examines Development permits, determines actions to take
func ExamineDevelopmentPermits() error {
	fetchedDevelopmentPermits, storedDevelopmentPermits, err := developmentpermit.GetDevelopmentPermits()
	if err != nil {
		return err
	}

	// Get development permit actions
	fileActions := developmentpermit.GetDevelopmentPermitActions(fetchedDevelopmentPermits, storedDevelopmentPermits)

	// Get Repository Id
	repositoryId, repoIdErr := githubdiscussions.GetRepository(config.Config.GithubDiscussions.Owner, config.Config.GithubDiscussions.Repository)
	if repoIdErr != nil {
		return repoIdErr
	}

	// Get github discussion catagories
	catagories, catagoriesErr := githubdiscussions.GetDiscussionCategories(config.Config.GithubDiscussions.Owner, config.Config.GithubDiscussions.Repository)
	if err != nil {
		return catagoriesErr
	}

	developmentPermitCatagory := githubdiscussions.FindCatagory(catagories, "Development Permits")
	if developmentPermitCatagory == nil {
		return fmt.Errorf("Error: Unable to idetify the Github Discussion Catagory ID for 'Development Permits'")
	}

	// Do actions for DPs
	for _, val := range fileActions {
		//Create
		if val.Action == "CREATE" {
			fmt.Printf("Development Permit %s:\n\tCreating Discussion...\n", val.PermitNum)

			//Find or create the discussion
			discussionId, err := FindOrCreateDiscussion(val.PermitNum, repositoryId, developmentPermitCatagory.ID, val.Message)
			if err != nil {
				fmt.Printf("\tFailed to create discussion. Error: %s\n", err.Error())
			}
			fmt.Printf("\tDiscussion Created!")

			//Append or Update change to stored DPs to be saved
			createdDP := developmentpermit.FindDevelopmentPermit(fetchedDevelopmentPermits, val.PermitNum)
			createdDP.GithubDiscussionId = &discussionId
			storedDevelopmentPermits = developmentpermit.UpsertDevelopmentPermit(storedDevelopmentPermits, *createdDP)
		}

		//Update
		if val.Action == "UPDATE" || val.Action == "CLOSE" {
			fmt.Printf("Development Permit %s:\n\tUpdating Discussion...\n", val.PermitNum)
			storedDP := developmentpermit.FindDevelopmentPermit(storedDevelopmentPermits, val.PermitNum)
			_, updateErr := githubdiscussions.AddDiscussionComment(*storedDP.GithubDiscussionId, val.Message)
			if updateErr != nil {
				fmt.Printf("\tFailed to comment on discussion. Error: %s\n", updateErr.Error())
				continue
			}

			//Append or Update change to stored DPs to be saved
			updatedDP := developmentpermit.FindDevelopmentPermit(fetchedDevelopmentPermits, val.PermitNum)
			updatedDP.GithubDiscussionId = storedDP.GithubDiscussionId
			updatedDP.GithubDiscussionClosed = storedDP.GithubDiscussionClosed

			if val.Action == "CLOSE" {
				fmt.Printf("\tClosing Discussion")

				closeErr := githubdiscussions.CloseDiscussion(*storedDP.GithubDiscussionId)
				if closeErr != nil {
					fmt.Printf("\tFailed to close discussion. Error: %s\n", closeErr.Error())
					continue
				}
				updatedDP.GithubDiscussionClosed = true
			}

			storedDevelopmentPermits = developmentpermit.UpsertDevelopmentPermit(storedDevelopmentPermits, *updatedDP)
		}

		break
	}

	// Save Development Permits
	developmentpermit.SaveDevelopmentPermits(storedDevelopmentPermits)

	return nil
}

// FindOrCreateDiscussion - attempts to find a discussion by its title, if not found creates it
func FindOrCreateDiscussion(permitNum, repositoryId, developmentPermitCategoryId, message string) (string, error) {
	discussionId, err := githubdiscussions.FindDiscussionByTitle(permitNum)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	// Only create a discussion if it does not already exist
	if discussionId == "" {
		discussionId, err = githubdiscussions.CreateDiscussion(repositoryId, developmentPermitCategoryId, permitNum, message)
		if err != nil {
			fmt.Println(err.Error())
			return "", err
		}
	}
	return discussionId, nil
}
