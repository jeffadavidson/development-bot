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

	for _, val := range fileActions {
		if val.Action == "CREATE" {
			fmt.Println("------------------------------------------------")
			fmt.Println("RepositoryId " + repositoryId)
			fmt.Println("Discussion catagopry " + developmentPermitCatagory.ID)
			fmt.Println("Will Create discussion for " + val.PermitNum)
			fmt.Println(val.Message)

			//Find or create the discussion
			discussionId, err := FindOrCreateDiscussion(val.PermitNum, repositoryId, developmentPermitCatagory.ID, val.Message)
			if err != nil {
				return err
			}

			fmt.Println("Discussion ID: " + discussionId)

			//Append change to stored DPs to be saved
			createdDP := developmentpermit.FindDevelopmentPermit(fetchedDevelopmentPermits, val.PermitNum)
			createdDP.GithubDiscussionId = discussionId
			storedDevelopmentPermits = append(storedDevelopmentPermits, *createdDP)
		}
	}

	// Save Development Permits
	developmentpermit.SaveDevelopmentPermits(storedDevelopmentPermits)

	// Send actions to GH actions
	// Update stored DP
	// Add action log for changes

	//TODO: What if update is not in the lookback window?

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
