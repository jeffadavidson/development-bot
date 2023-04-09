package examinedata

import (
	//"fmt"

	"fmt"

	"github.com/jeffadavidson/development-bot/interactions/githubdiscussions"
	"github.com/jeffadavidson/development-bot/objects/developmentpermit"
	"github.com/jeffadavidson/development-bot/utilities/config"
)

var repositoryId string
var repositoryCatagories []githubdiscussions.GithubDiscussionCatagory

func ManualInit() error {
	// Get Repository Id
	repoId, repoIdErr := githubdiscussions.GetRepository(config.Config.GithubDiscussions.Owner, config.Config.GithubDiscussions.Repository)
	if repoIdErr != nil {
		return repoIdErr
	}

	// Get github discussion catagories
	catagories, catagoriesErr := githubdiscussions.GetDiscussionCategories(config.Config.GithubDiscussions.Owner, config.Config.GithubDiscussions.Repository)
	if catagoriesErr != nil {
		return catagoriesErr
	}

	repositoryId = repoId
	repositoryCatagories = catagories

	return nil
}

// DevelopmentPermits - Evaluates Development permits, determines actions to take
func DevelopmentPermits() error {
	// Get the catagoryId for 'development permits'
	developmentPermitCatagory := githubdiscussions.FindCatagory(repositoryCatagories, "Development Permits")
	if developmentPermitCatagory == nil {
		return fmt.Errorf("Error: Unable to idetify the Github Discussion Catagory ID for 'Development Permits'")
	}

	dpErr := developmentpermit.EvaluateDevelopmentPermits(repositoryId, developmentPermitCatagory.ID)
	if dpErr != nil {
		return dpErr
	}

	return nil
}
