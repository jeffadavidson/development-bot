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

func Test_ArePointersEqual_BothNotNil(t *testing.T) {
	s1 := "test-string"
	p1 := &s1
	p2 := &s1
	assert.True(t, ArePointersEqual(p1, p2), "Expected ArePointersEqual(%v, %v) to be true, but got false", p1, p2)
}

func Test_ArePointersEqual_BothNil(t *testing.T) {
	var p3 *string
	var p4 *string
	assert.True(t, ArePointersEqual(p3, p4), "Expected ArePointersEqual(%v, %v) to be true, but got false", p3, p4)
}

func Test_ArePointersEqual_OneNil(t *testing.T) {
	s2 := "test-string"
	p5 := &s2
	var p6 *string
	assert.False(t, ArePointersEqual(p5, p6), "Expected ArePointersEqual(%v, %v) to be false, but got true", p5, p6)
}
