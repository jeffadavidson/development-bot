package toolbox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains_ElementInSlice_String(t *testing.T) {
	mySlice := []string{"apple", "banana", "cherry"}
	result := SliceContains(mySlice, "banana")
	assert.Equal(t, true, result)
}

func TestContains_ElementInSlice_Int(t *testing.T) {
	mySlice := []int{1, 2, 3, 4, 5}
	result := SliceContains(mySlice, 1)
	assert.Equal(t, true, result)
}

func TestContains_ElementNotInSlice(t *testing.T) {
	mySlice := []int{1, 2, 3, 4, 5}
	result := SliceContains(mySlice, 6)
	assert.Equal(t, false, result)
}

func TestContains_EmptySlice(t *testing.T) {
	var mySlice []float64
	result := SliceContains(mySlice, 1.2)
	assert.Equal(t, false, result)
}
