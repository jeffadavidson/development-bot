package fileio

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileContents(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}
	defer os.Remove(tempFile.Name())

	content := "example content"
	if _, err := tempFile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write content to temporary file: %s", err)
	}

	if err := tempFile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %s", err)
	}

	fileBytes, err := GetFileContents(tempFile.Name())
	assert.NoError(t, err)
	assert.Equal(t, string(fileBytes), content)
}

func Test_WriteFileContents(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	testData := []byte("test data")

	err = WriteFileContents(tmpfile.Name(), testData)
	assert.NoError(t, err)

	contents, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(testData), string(contents))
}
