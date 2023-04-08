package developmentpermit

import (
	"fmt"
	"net/url"

	"golang.org/x/exp/slices"
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

// CreateInformationMessage - Builds an information message from the development permit
func (dp DevelopmentPermit) CreateInformationMessage() string {
	var lineSeparator string = "\n"
	var message string

	message += fmt.Sprintf("Development Permit Application: %v", dp.PermitNum)
	if dp.AppliedDate != "" {
		message += fmt.Sprintf("%vDate Applied: %v", lineSeparator, dp.AppliedDate)
	}
	if dp.Address != "" {
		message += fmt.Sprintf("%vAddress: %v", lineSeparator, dp.Address)
	}
	if dp.CommunityName != "" {
		message += fmt.Sprintf("%vCommunity: %v", lineSeparator, dp.CommunityName)
	}
	if dp.Applicant != "" {
		message += fmt.Sprintf("%vApplicant: %v", lineSeparator, dp.Applicant)
	}
	if dp.Description != "" {
		message += fmt.Sprintf("%vDescription: %v", lineSeparator, dp.Description)
	}
	if dp.LandUseDistrictDesc != "" {
		message += fmt.Sprintf("%vCurrent Land Use: %v", lineSeparator, dp.LandUseDistrictDesc)
	}
	if dp.StatusCurrent != "" {
		message += fmt.Sprintf("%vCurrent Permit Status: %v", lineSeparator, dp.StatusCurrent)
	}
	if dp.MustCommenceDate != "" {
		message += fmt.Sprintf("%vMust Comment By Date: %v", lineSeparator, dp.MustCommenceDate)
	}
	if dp.PermittedDiscretion != "" {
		message += fmt.Sprintf("%vApplication Type: %v", lineSeparator, dp.PermittedDiscretion)
	}
	if dp.Decision != "" {
		message += fmt.Sprintf("%vDecision: %v", lineSeparator, dp.Decision)
	}
	if dp.ReleaseDate != "" {
		message += fmt.Sprintf("%vRelease Date: %v", lineSeparator, dp.ReleaseDate)
	}

	// Add google and dmap links
	message += lineSeparator + lineSeparator
	message += fmt.Sprintf("%vDmap: https://dmap.calgary.ca/?find=%v", lineSeparator, dp.PermitNum)
	if dp.Address != "" {
		message += fmt.Sprintf("%vGoogle Maps: https://maps.google.com/?q=%v", lineSeparator, url.QueryEscape(fmt.Sprintf("%v, Calgary, Alberta", dp.Address)))
	}
	message += lineSeparator

	return message
}

// FindDevelopmentPermit - finds a development permit in a list of permits
func FindDevelopmentPermit(searchSlice []DevelopmentPermit, permitNum string) *DevelopmentPermit {
	foundIndex := slices.IndexFunc(searchSlice, func(c DevelopmentPermit) bool { return c.PermitNum == permitNum })
	if foundIndex == -1 {
		return nil
	}

	return &searchSlice[foundIndex]
}
