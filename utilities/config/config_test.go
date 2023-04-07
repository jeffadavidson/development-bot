package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseConfig_NoBoundingBox(t *testing.T) {
	configYaml := []byte(`
  runmode: DEVELOPMENT
  neighborhood:
    name: Killarney
`)

	err := parseConfig(configYaml)
	assert.Equal(t, nil, err)
}

func Test_ParseConfig_ValidConfig(t *testing.T) {
	configYaml := []byte(`
  runmode: DEVELOPMENT
  neighborhood:
    name: Killarney
    bounding-box:
      north-latitude: 51.038912
      east-longitude: -114.117927
      south-latitude: 51.022361
      west-longitude: -114.142638 
`)

	err := parseConfig(configYaml)
	assert.Equal(t, nil, err)
}

func Test_PaseConfig_InvalidLatitude(t *testing.T) {
	configYaml := []byte(`
  runmode: DEVELOPMENT
  neighborhood:
    name: Killarney
    bounding-box:
      north-latitude: invalid
      east-longitude: -114.117927
      south-latitude: 51.022361
      west-longitude: -114.142638 
`)

	err := parseConfig(configYaml)
	assert.NotEqual(t, nil, err)
}
