package calgaryopendata

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jeffadavidson/development-bot/utilities/config"
	"github.com/jeffadavidson/development-bot/utilities/simplehttp"
)

type DevelopmentPermit struct {
	Point                  Point  `json:"point"`
	PermitNum              string `json:"permitnum"`
	Address                string `json:"address"`
	Applicant              string `json:"applicant"`
	Category               string `json:"category"`
	Description            string `json:"description"`
	ProposedUseCode        string `json:"proposedusecode"`
	ProposedUseDesc        string `json:"proposedusedescription"`
	PermittedDiscretion    string `json:"permitteddiscretionary"`
	LandUseDistrict        string `json:"landusedistrict"`
	LandUseDistrictDesc    string `json:"landusedistrictdescription"`
	StatusCurrent          string `json:"statuscurrent"`
	AppliedDate            string `json:"applieddate"`
	DecisionDate           string `json:"decisiondate"`
	ReleaseDate            string `json:"releasedate"`
	MustCommenceDate       string `json:"mustcommencedate"`
	Decision               string `json:"decision"`
	DecisionBy             string `json:"decisionby"`
	CommunityCode          string `json:"communitycode"`
	CommunityName          string `json:"communityname"`
	Ward                   string `json:"ward"`
	Quadrant               string `json:"quadrant"`
	Latitude               string `json:"latitude"`
	Longitude              string `json:"longitude"`
	LocationCount          string `json:"locationcount"`
	LocationTypes          string `json:"locationtypes"`
	LocationAddresses      string `json:"locationaddresses"`
	LocationsGeoJSON       string `json:"locationsgeojson"`
	LocationsWKT           string `json:"locationswkt"`
	GithubDiscussionClosed bool   `json:"github_discussion_closed"`
}

type Point struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

func GetDevelopmentPermits() ([]DevelopmentPermit, error) {
	var developmentPermits []DevelopmentPermit

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

func ParseDevelopmentPermits(developmentPermitByte []byte) ([]DevelopmentPermit, error) {
	var developmentPermits []DevelopmentPermit
	err := json.Unmarshal(developmentPermitByte, &developmentPermits)
	if err != nil {
		return developmentPermits, fmt.Errorf("failed to parse development permit json. Error: %s", err.Error())
	}

	return developmentPermits, nil
}
