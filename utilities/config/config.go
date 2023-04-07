package config

import (
	"fmt"

	"github.com/jeffadavidson/development-bot/utilities/exit"
	"github.com/jeffadavidson/development-bot/utilities/fileio"

	"gopkg.in/yaml.v3"
)

var configFilePath string = "config.yaml"

type DevBot struct {
	RunMode      string       `yaml:"runmode"`
	Neighborhood Neighborhood `yaml:"neighborhood"`
}

type Neighborhood struct {
	Name        string      `yaml:"name"`
	BoundingBox BoundingBox `yaml:"bounding-box"`
}

type BoundingBox struct {
	NorthLatitude float64 `yaml:"north-latitude"`
	EastLongitude float64 `yaml:"east-longitude"`
	SouthLatitude float64 `yaml:"south-latitude"`
	WestLongitude float64 `yaml:"west-longitude"`
}

var Config DevBot

func ManualInit() {
	loaderr := loadConfig(configFilePath)
	if loaderr != nil {
		exit.ExitError(loaderr)
	}
}

func loadConfig(filePath string) error {
	// Load the YAML file into a byte slice
	yamlFile, err := fileio.GetFileContents("config.yaml")
	if err != nil {
		return fmt.Errorf("error reading YAML file: %v", err)
	}

	parseerr := parseConfig(yamlFile)
	if parseerr != nil {
		return parseerr
	}

	return nil
}

func parseConfig(configBytes []byte) error {
	// Unmarshal the YAML data into a DevBot struct
	if err := yaml.Unmarshal(configBytes, &Config); err != nil {
		return fmt.Errorf("error unmarshalling YAML data: %v", err)
	}

	return nil
}
