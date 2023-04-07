package config

import (
	"fmt"
	"io/ioutil"

	"github.com/jeffadavidson/development-bot/utilities/exit"
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
	NWLatitude  float64 `yaml:"nw-latitude"`
	NWLongitude float64 `yaml:"nw-longitude"`
	NELatitude  float64 `yaml:"ne-latitude"`
	NELongitude float64 `yaml:"ne-longitude"`
	SWLatitude  float64 `yaml:"sw-latitude"`
	SWLongitude float64 `yaml:"sw-longitude"`
	SELatitude  float64 `yaml:"se-latitude"`
	SELongitude float64 `yaml:"se-longitude"`
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
	yamlFile, err := ioutil.ReadFile("config.yaml")
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
