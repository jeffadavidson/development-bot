package config

import (
	"fmt"
	"os"

	"github.com/jeffadavidson/development-bot/utilities/fileio"

	"gopkg.in/yaml.v3"
)

var configFilePath string = "config.yaml"

type DevBot struct {
	RunMode           string           `yaml:"runmode"`
	GithubDiscussions GithubDiscussion `yaml:"github-discussion"`
	Neighborhood      Neighborhood     `yaml:"neighborhood"`
}

type GithubDiscussion struct {
	Owner      string `yaml:"owner"`
	Repository string `yaml:"repository"`
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

func ManualInit() error {
	loaderr := loadConfig(configFilePath)
	if loaderr != nil {
		return loaderr
	}

	// Check for environmental variable overrides
	runmode := os.Getenv("DEVELOPMENT_BOT_RUNMODE")
	if runmode != "" {
		Config.RunMode = runmode
	}

	fmt.Println("CONFIG")
	fmt.Println(Config.RunMode)

	fmt.Println("CONFIG")
	return nil
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
