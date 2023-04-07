package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	RunMode      string             `yaml:"runmode"`
	Neighborhood NeighborhoodConfig `yaml:"neighborhood"`
}

type NeighborhoodConfig struct {
	BoundingBox BoundingBoxConfig `yaml:"bounding-box"`
}

type BoundingBoxConfig struct {
	NWLatitude  float64 `yaml:"nw-latitude"`
	NWLongitude float64 `yaml:"nw-longitude"`
	NELatitude  float64 `yaml:"ne-latitude"`
	NELongitude float64 `yaml:"ne-longitude"`
	SWLatitude  float64 `yaml:"sw-latitude"`
	SWLongitude float64 `yaml:"sw-longitude"`
	SELatitude  float64 `yaml:"se-latitude"`
	SELongitude float64 `yaml:"se-longitude"`
}

var theconfig Config

func init() {
	fmt.Println("config init")

	// Load the YAML file into a byte slice
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Printf("error reading YAML file: %v", err)
	}

	// Unmarshal the YAML data into a Config struct

	if err := yaml.Unmarshal(yamlFile, &theconfig); err != nil {
		fmt.Printf("error unmarshalling YAML data: %v", err)
	}
}
