package fileio

import (
	"fmt"
	"io/ioutil"
)

func GetFileContents(filepath string) ([]byte, error) {
	var fileBytes []byte
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		// Handle error
		return fileBytes, fmt.Errorf("Error reading file. Error: %s", err.Error())
	}
	fileBytes = bytes

	return fileBytes, nil
}

func WriteFileContents(filepath string, fileBytes []byte) error {
	err := ioutil.WriteFile(filepath, fileBytes, 0644)
	if err != nil {
		return fmt.Errorf("Error writing file. Error: %s", err.Error())
	}

	return nil
}
