package calgaryopendata

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jeffadavidson/development-bot/utilities/config"
	"github.com/jeffadavidson/development-bot/utilities/simplehttp"
)

func GetDevelopmentPermits() ([]byte, error) {
	var developmentPermits []byte

	// Build URL
	baseUrl := "https://data.calgary.ca/resource/6933-unw5.json"
	// Only look back 3 months for recent activity
	threeMonthsAgo := time.Now().AddDate(0, -3, 0).Format("2006-01-02T15:04:05.000")
	query := fmt.Sprintf("$query=SELECT * WHERE applieddate > '%s' AND latitude BETWEEN '%f' AND '%f' AND longitude BETWEEN '%f' AND '%f' ORDER BY applieddate DESC", threeMonthsAgo, config.Config.Neighborhood.BoundingBox.SouthLatitude, config.Config.Neighborhood.BoundingBox.NorthLatitude, config.Config.Neighborhood.BoundingBox.EastLongitude, config.Config.Neighborhood.BoundingBox.WestLongitude)
	url := fmt.Sprintf("%s?%s", baseUrl, query)

	response, err := simplehttp.SimpleGet(url, make(map[string]string))
	if err != nil {
		return developmentPermits, fmt.Errorf("error getting development permits from Calgary Open Data. Error: %s", err.Error())
	}
	if response.StatusCode != http.StatusOK {
		return developmentPermits, fmt.Errorf("error getting development permits from Calgary Open Data. Http Status %d", response.StatusCode)
	}

	developmentPermits = response.Body

	return developmentPermits, nil
}

func GetRezoningApplications() ([]byte, error) {
	var developmentPermits []byte

	// Build URL
	baseUrl := "https://data.calgary.ca/resource/33vi-ew4s.json"
	// Only look back 3 months for recent activity
	threeMonthsAgo := time.Now().AddDate(0, -3, 0).Format("2006-01-02T15:04:05.000")
	query := fmt.Sprintf("$query=SELECT * WHERE applieddate > '%s' AND latitude BETWEEN '%f' AND '%f' AND longitude BETWEEN '%f' AND '%f' ORDER BY applieddate DESC", threeMonthsAgo, config.Config.Neighborhood.BoundingBox.SouthLatitude, config.Config.Neighborhood.BoundingBox.NorthLatitude, config.Config.Neighborhood.BoundingBox.EastLongitude, config.Config.Neighborhood.BoundingBox.WestLongitude)
	url := fmt.Sprintf("%s?%s", baseUrl, query)

	response, err := simplehttp.SimpleGet(url, make(map[string]string))
	if err != nil {
		return developmentPermits, fmt.Errorf("error getting development permits from Calgary Open Data. Error: %s", err.Error())
	}
	if response.StatusCode != http.StatusOK {
		return developmentPermits, fmt.Errorf("error getting development permits from Calgary Open Data. Http Status %d", response.StatusCode)
	}

	developmentPermits = response.Body

	return developmentPermits, nil
}
