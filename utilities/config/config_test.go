package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseConfig_NoBoundingBox(t *testing.T) {
	configYaml := []byte(`
runmode: DEVELOPMENT
neighborhood:
  name: TestNeighborhood
`)

	err := parseConfig(configYaml)
	assert.Equal(t, nil, err)
}

func Test_ParseConfig_ValidConfig(t *testing.T) {
	configYaml := []byte(`
runmode: DEVELOPMENT
neighborhood:
  name: TestNeighborhood
  bounding-box:
    nw-latitude: 51.038912
    nw-longitude: -114.142638
    ne-latitude: 51.038912
    ne-longitude: -114.117927
    sw-latitude: 51.022361
    sw-longitude: -114.142638
    se-latitude: 51.022361
    se-longitude: -114.117927
`)

	err := parseConfig(configYaml)
	assert.Equal(t, nil, err)
}

func Test_PaseConfig_InvalidLatitude(t *testing.T) {
	configYaml := []byte(`
runmode: DEVELOPMENT
neighborhood:
  name: TestNeighborhood
  bounding-box:
    nw-latitude: invalid
    nw-longitude: -114.142638
    ne-latitude: 51.038912
    ne-longitude: -114.117927
    sw-latitude: 51.022361
    sw-longitude: -114.142638
    se-latitude: 51.022361
    se-longitude: -114.117927
`)

	err := parseConfig(configYaml)
	assert.NotEqual(t, nil, err)
}
