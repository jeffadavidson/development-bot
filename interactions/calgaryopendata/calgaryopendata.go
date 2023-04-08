package calgaryopendata

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jeffadavidson/development-bot/objects/developmentpermit"
	"github.com/jeffadavidson/development-bot/utilities/config"
	"github.com/jeffadavidson/development-bot/utilities/simplehttp"
)

func GetDevelopmentPermits() ([]developmentpermit.DevelopmentPermit, error) {
	var developmentPermits []developmentpermit.DevelopmentPermit

	// Build URL
	baseUrl := "https://data.calgary.ca/resource/6933-unw5.json"
	query := fmt.Sprintf("$query=SELECT * WHERE applieddate > '2022-01-01T00:00:00.000' AND latitude BETWEEN '%f' AND '%f' AND longitude BETWEEN '%f' AND '%f' ORDER BY applieddate DESC", config.Config.Neighborhood.BoundingBox.SouthLatitude, config.Config.Neighborhood.BoundingBox.NorthLatitude, config.Config.Neighborhood.BoundingBox.EastLongitude, config.Config.Neighborhood.BoundingBox.WestLongitude)
	url := fmt.Sprintf("%s?%s", baseUrl, query)

	response, err := simplehttp.SimpleGet(url, make(map[string]string))
	if err != nil {
		return developmentPermits, fmt.Errorf("error getting development permits from Calgary Open Data. Error: %s", err.Error())
	}
	if response.StatusCode != http.StatusOK {
		return developmentPermits, fmt.Errorf("error getting development permits from Calgary Open Data. Http Status %d", response.StatusCode)
	}

	var parseErr error
	developmentPermits, parseErr = ParseDevelopmentPermits(response.Body)
	if parseErr != nil {
		return developmentPermits, fmt.Errorf("error parsing development permits. Error %s", parseErr.Error())
	}

	return developmentPermits, nil
}

func ParseDevelopmentPermits(developmentPermitByte []byte) ([]developmentpermit.DevelopmentPermit, error) {
	var developmentPermits []developmentpermit.DevelopmentPermit
	err := json.Unmarshal(developmentPermitByte, &developmentPermits)
	if err != nil {
		return developmentPermits, fmt.Errorf("failed to parse development permit json. Error: %s", err.Error())
	}

	return developmentPermits, nil
}
