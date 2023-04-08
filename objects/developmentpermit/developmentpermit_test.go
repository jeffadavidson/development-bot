package developmentpermit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FindDevelopmentPermit_ItemFound(t *testing.T) {
	dp1 := DevelopmentPermit{PermitNum: "123"}
	dp2 := DevelopmentPermit{PermitNum: "456"}
	dp3 := DevelopmentPermit{PermitNum: "789"}
	searchSlice := []DevelopmentPermit{dp1, dp2, dp3}

	result := FindDevelopmentPermit(searchSlice, "456")

	assert.NotNil(t, result)
	assert.Equal(t, &dp2, result)
}

func Test_FindDevelopmentPermit_ItemNotFound(t *testing.T) {
	dp1 := DevelopmentPermit{PermitNum: "123"}
	//	dp2 := DevelopmentPermit{PermitNum: "456"}
	dp3 := DevelopmentPermit{PermitNum: "789"}
	searchSlice := []DevelopmentPermit{dp1, dp3}

	result := FindDevelopmentPermit(searchSlice, "456")

	assert.Nil(t, result)
}

func Test_FindDevelopmentPermit_EmptySlice(t *testing.T) {
	searchSlice := []DevelopmentPermit{}

	result := FindDevelopmentPermit(searchSlice, "456")

	assert.Nil(t, result)
}
