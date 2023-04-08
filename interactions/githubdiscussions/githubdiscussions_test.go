package githubdiscussions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FindDiscussionCatagory_ItemFound(t *testing.T) {
	dp1 := GithubDiscussionCatagory{Name: "123"}
	dp2 := GithubDiscussionCatagory{Name: "456"}
	dp3 := GithubDiscussionCatagory{Name: "789"}
	searchSlice := []GithubDiscussionCatagory{dp1, dp2, dp3}

	result := FindCatagory(searchSlice, "456")

	assert.NotNil(t, result)
	assert.Equal(t, &dp2, result)
}

func Test_FindDiscussionCatagory_ItemNotFound(t *testing.T) {
	dp1 := GithubDiscussionCatagory{Name: "123"}
	//	dp2 := GithubDiscussionCatagory{Name: "456"}
	dp3 := GithubDiscussionCatagory{Name: "789"}
	searchSlice := []GithubDiscussionCatagory{dp1, dp3}

	result := FindCatagory(searchSlice, "456")

	assert.Nil(t, result)
}

func Test_FindDiscussionCatagory_EmptySlice(t *testing.T) {
	searchSlice := []GithubDiscussionCatagory{}

	result := FindCatagory(searchSlice, "456")

	assert.Nil(t, result)
}
