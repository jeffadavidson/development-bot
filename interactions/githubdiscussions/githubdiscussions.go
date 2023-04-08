package githubdiscussions

import (
	"context"
	"fmt"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"
)

type GithubDiscussionCatagory struct {
	ID          string
	Name        string
	Description string
}

var githubClient *githubv4.Client

// init - Initilize client on startup
func init() {
	// Get PAT from environment variable
	githubPat := os.Getenv("GITHUB_PAT")
	if githubPat == "" {
		panic("error, github pat not set")
	}

	githubClient = newGitHubClient(githubPat)
}

// newGitHubClient - creates a new github client
func newGitHubClient(token string) *githubv4.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	return client
}

// GetRepository - Gets the repository id for a repository
func GetRepository(owner string, repositoryName string) (string, error) {
	var query struct {
		Repository struct {
			ID string
		} `graphql:"repository(owner: $owner, name: $name)"`
	}

	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(repositoryName),
	}

	err := githubClient.Query(context.Background(), &query, variables)
	if err != nil {
		return "", fmt.Errorf("failed to get repository ID: %w", err)
	}

	return string(query.Repository.ID), nil
}

// GetDiscussionCategories - Gets the discussion catagories for an owner and repo
func GetDiscussionCategories(owner string, repository string) ([]GithubDiscussionCatagory, error) {
	var query struct {
		Repository struct {
			DiscussionCategories struct {
				Nodes []struct {
					Id          string
					Name        string
					Description string
				}
			} `graphql:"discussionCategories(first: 100)"`
		} `graphql:"repository(owner:$owner, name:$repository)"`
	}

	variables := map[string]interface{}{
		"owner":      githubv4.String(owner),
		"repository": githubv4.String(repository),
	}

	err := githubClient.Query(context.Background(), &query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to get discussion categories: %w", err)
	}

	categories := []GithubDiscussionCatagory{}
	for _, node := range query.Repository.DiscussionCategories.Nodes {
		categories = append(categories, GithubDiscussionCatagory{ID: node.Id, Name: node.Name, Description: node.Description})
	}

	return categories, nil
}

// FindCatagory - Finds a catagory by name in a slice of catagories
func FindCatagory(catagories []GithubDiscussionCatagory, catagoryName string) *GithubDiscussionCatagory {
	foundIndex := slices.IndexFunc(catagories, func(c GithubDiscussionCatagory) bool { return c.Name == catagoryName })
	if foundIndex == -1 {
		return nil
	}

	return &catagories[foundIndex]
}

func FindDiscussionByTitle(title string) (string, error) {
	var q struct {
		Search struct {
			Edges []struct {
				Node struct {
					Discussion struct {
						ID    string
						Title string
					} `graphql:"... on Discussion"`
				}
			}
		} `graphql:"search(query: $query, type: DISCUSSION, first: 1)"`
	}
	query := fmt.Sprintf("in:title %s", title)

	variables := map[string]interface{}{
		"query": githubv4.String(query),
	}
	err := githubClient.Query(context.Background(), &q, variables)
	if err != nil {
		return "", err
	}
	if len(q.Search.Edges) == 0 {
		return "", fmt.Errorf("discussion not found")
	}
	return q.Search.Edges[0].Node.Discussion.ID, nil
}

func CreateDiscussion(repositoryID, categoryID, title, body string) (string, error) {
	var mutation struct {
		CreateDiscussion struct {
			Discussion struct {
				ID string
			}
		} `graphql:"createDiscussion(input: $input)"`
	}

	input := githubv4.CreateDiscussionInput{
		RepositoryID: githubv4.String(repositoryID),
		CategoryID:   githubv4.String(categoryID),
		Title:        githubv4.String(title),
		Body:         githubv4.String(body),
	}

	err := githubClient.Mutate(context.Background(), &mutation, input, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create discussion: %w", err)
	}

	return mutation.CreateDiscussion.Discussion.ID, nil
}
