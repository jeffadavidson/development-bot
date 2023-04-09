package developmentpermit

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/jeffadavidson/development-bot/interactions/calgaryopendata"
	"github.com/jeffadavidson/development-bot/objects/fileaction"
	"github.com/jeffadavidson/development-bot/utilities/fileio"
	"github.com/jeffadavidson/development-bot/utilities/toolbox"
	"golang.org/x/exp/slices"
)

type DevelopmentPermit struct {
	Point                  Point   `json:"point"`
	PermitNum              string  `json:"permitnum"`
	Address                *string `json:"address"`
	Applicant              *string `json:"applicant"`
	Category               *string `json:"category"`
	Description            *string `json:"description"`
	ProposedUseCode        *string `json:"proposedusecode"`
	ProposedUseDesc        *string `json:"proposedusedescription"`
	PermittedDiscretion    *string `json:"permitteddiscretionary"`
	LandUseDistrict        *string `json:"landusedistrict"`
	LandUseDistrictDesc    *string `json:"landusedistrictdescription"`
	StatusCurrent          string  `json:"statuscurrent"`
	AppliedDate            *string `json:"applieddate"`
	CommunityCode          *string `json:"communitycode"`
	CommunityName          *string `json:"communityname"`
	Ward                   *string `json:"ward"`
	Quadrant               *string `json:"quadrant"`
	Latitude               *string `json:"latitude"`
	Longitude              *string `json:"longitude"`
	LocationCount          *string `json:"locationcount"`
	LocationTypes          *string `json:"locationtypes"`
	LocationAddresses      *string `json:"locationaddresses"`
	LocationsGeoJSON       *string `json:"locationsgeojson"`
	LocationsWKT           *string `json:"locationswkt"`
	DecisionDate           *string `json:"decisiondate"`
	MustCommenceDate       *string `json:"mustcommencedate"`
	Decision               *string `json:"decision"`
	DecisionBy             *string `json:"decisionby"`
	ReleaseDate            *string `json:"releasedate"`
	GithubDiscussionId     *string `json:"github_discussion_id"`
	GithubDiscussionClosed bool    `json:"github_discussion_closed"`
}

type Point struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// CreateInformationMessage - Builds an information message from the development permit
func (dp DevelopmentPermit) CreateInformationMessage() string {
	var lineSeparator string = "\n"
	var message string = ""

	// Markdown
	message += "## About\n\n"
	message += fmt.Sprintf("**Permit Number:** %v", dp.PermitNum)
	if dp.AppliedDate != nil {
		appliedDate, err := time.Parse("2006-01-02T15:04:05.000", *dp.AppliedDate)
		if err != nil {
			message += fmt.Sprintf("%v**Date Applied:** %v", lineSeparator, *dp.AppliedDate)
		} else {
			dateStr := appliedDate.Format("2006-01-02")
			message += fmt.Sprintf("%v**Date Applied:** %v", lineSeparator, dateStr)
		}
	}
	if dp.Address != nil {
		message += fmt.Sprintf("%v**Address:** %v", lineSeparator, *dp.Address)
	}
	if dp.CommunityName != nil {
		message += fmt.Sprintf("%v**Community:** %v", lineSeparator, *dp.CommunityName)
	}
	if dp.Applicant != nil {
		message += fmt.Sprintf("%v**Applicant:** %v", lineSeparator, *dp.Applicant)
	}
	if dp.Description != nil {
		message += fmt.Sprintf("%v**Description:** %v", lineSeparator, *dp.Description)
	}
	if dp.LandUseDistrictDesc != nil {
		message += fmt.Sprintf("%v**Current Land Use:** %v", lineSeparator, *dp.LandUseDistrictDesc)
	}
	if dp.StatusCurrent != "" {
		message += fmt.Sprintf("%v**Permit Status:** %v", lineSeparator, dp.StatusCurrent)
	}
	if dp.MustCommenceDate != nil {
		message += fmt.Sprintf("%v**Must Comment By Date:** %v", lineSeparator, *dp.MustCommenceDate)
	}
	if dp.PermittedDiscretion != nil {
		message += fmt.Sprintf("%v**Application Type:** %v", lineSeparator, *dp.PermittedDiscretion)
	}
	if dp.Decision != nil {
		message += fmt.Sprintf("%v**Decision:** %v", lineSeparator, *dp.Decision)
	}
	if dp.ReleaseDate != nil {
		message += fmt.Sprintf("%v**Release Date:** %v", lineSeparator, *dp.ReleaseDate)
	}

	// Add google and dmap links
	message += "\n## Links\n\n"
	message += lineSeparator + lineSeparator
	message += fmt.Sprintf("%v[Development Map](https://dmap.calgary.ca/?find=%v)", lineSeparator, dp.PermitNum)
	if dp.Address != nil {
		message += fmt.Sprintf("%v [Google Maps](https://maps.google.com/?q=%v)", lineSeparator, url.QueryEscape(fmt.Sprintf("%v, Calgary, Alberta", *dp.Address)))
	}
	message += lineSeparator

	return message
}

// getDevelopmentPermits - Gets fetched development permits, gets stored development permits
func GetDevelopmentPermits() ([]DevelopmentPermit, []DevelopmentPermit, error) {
	// Load existing development permits
	storedDevelopmentPermitsBytes, loadErr := fileio.GetFileContents("./data/development-permits.json")
	if loadErr != nil {
		return nil, nil, loadErr
	}
	storedDevelopmentPermits, parseErr := parseDevelopmentPermits(storedDevelopmentPermitsBytes)
	if parseErr != nil {
		return nil, nil, parseErr
	}

	//Get development Permits from calgary open data
	fetchedDevelopmentPermitsRaw, fetchErr := calgaryopendata.GetDevelopmentPermits()
	if fetchErr != nil {
		return nil, nil, fetchErr
	}
	fetchedDevelopmentPermits, parseErr2 := parseDevelopmentPermits(fetchedDevelopmentPermitsRaw)
	if parseErr2 != nil {
		return nil, nil, parseErr
	}

	return fetchedDevelopmentPermits, storedDevelopmentPermits, nil
}

func SaveDevelopmentPermits(permits []DevelopmentPermit) error {
	// Encode permits as JSON
	permitsBytes, encodeErr := json.MarshalIndent(permits, "", "  ")
	if encodeErr != nil {
		return encodeErr
	}

	writeErr := fileio.WriteFileContents("./data/development-permits.json", permitsBytes)
	if writeErr != nil {
		return writeErr
	}
	return nil
}

func parseDevelopmentPermits(developmentPermitByte []byte) ([]DevelopmentPermit, error) {
	var developmentPermits []DevelopmentPermit
	err := json.Unmarshal(developmentPermitByte, &developmentPermits)
	if err != nil {
		return developmentPermits, fmt.Errorf("failed to parse development permit json. Error: %s", err.Error())
	}

	return developmentPermits, nil
}

// FindDevelopmentPermit - finds a development permit in a list of permits
func FindDevelopmentPermit(searchSlice []DevelopmentPermit, permitNum string) *DevelopmentPermit {
	foundIndex := slices.IndexFunc(searchSlice, func(c DevelopmentPermit) bool { return c.PermitNum == permitNum })
	if foundIndex == -1 {
		return nil
	}

	return &searchSlice[foundIndex]
}

// getDevelopmentPermitActions - For a list of fetched and stored development permits compares permits and gets a list of actions to execute
func GetDevelopmentPermitActions(fetchedDevelopmentPermits []DevelopmentPermit, storedDevelopmentPermits []DevelopmentPermit) []fileaction.FileAction {
	var fileActions []fileaction.FileAction
	for _, fetchedDP := range fetchedDevelopmentPermits {
		storedDpPtr := FindDevelopmentPermit(storedDevelopmentPermits, fetchedDP.PermitNum)
		if storedDpPtr == nil {
			fileActions = append(fileActions, fileaction.FileAction{PermitNum: fetchedDP.PermitNum, Action: "CREATE", Message: fetchedDP.CreateInformationMessage()})
		} else {
			storedDP := *storedDpPtr

			// Skip if discussion closed
			if storedDP.GithubDiscussionClosed {
				fileActions = append(fileActions, fileaction.FileAction{PermitNum: fetchedDP.PermitNum, Action: "SKIP"})
				continue
			}

			// Create if discussion does not exist but its still stored(from when we switched from slack)
			if storedDP.GithubDiscussionId == nil {
				fileActions = append(fileActions, fileaction.FileAction{PermitNum: fetchedDP.PermitNum, Action: "CREATE", Message: fetchedDP.CreateInformationMessage()})
				continue
			}

			hasUpdate, updateMessage := getDevelopmentPermitUpdates(fetchedDP, storedDP)
			toClose, closeMessage := isDevelopmentPermitClosed(fetchedDP, storedDP)
			var message string

			if hasUpdate && !toClose {
				message += updateMessage
				fileActions = append(fileActions, fileaction.FileAction{PermitNum: fetchedDP.PermitNum, Action: "UPDATE", Message: message})
			}
			if hasUpdate && toClose {
				message += updateMessage
				message += "\n"
				message += closeMessage

				fileActions = append(fileActions, fileaction.FileAction{PermitNum: fetchedDP.PermitNum, Action: "CLOSE", Message: message})
			}
			if !hasUpdate && toClose {
				message += closeMessage
				fileActions = append(fileActions, fileaction.FileAction{PermitNum: fetchedDP.PermitNum, Action: "CLOSE", Message: closeMessage})
			}
		}
	}

	return fileActions
}

// isDevelopmentPermitClosed - Checks if a development permit is ready to be closed
func isDevelopmentPermitClosed(fetchedDP DevelopmentPermit, storedDP DevelopmentPermit) (bool, string) {
	toClose := false
	closeMessage := ""

	// Check for close
	close_statuses := [3]string{"Released", "Cancelled", "Cancelled - Pending Refund"}
	if toolbox.SliceContains([]string(close_statuses[:]), fetchedDP.StatusCurrent) {
		toClose = true
		closeMessage = fmt.Sprintf("Closing file as it is in status '%s'", fetchedDP.StatusCurrent)
	}

	return toClose, closeMessage
}

// getDevelopmentPermitUpdates - Checks if a development permit needs updates
func getDevelopmentPermitUpdates(fetchedDP DevelopmentPermit, storedDP DevelopmentPermit) (bool, string) {
	hasUpdate := false
	updateMessage := ""

	// check status
	if fetchedDP.StatusCurrent != storedDP.StatusCurrent {
		hasUpdate = true
		updateMessage += fmt.Sprintf("Status updated from '%s' to '%s'\n", storedDP.StatusCurrent, fetchedDP.StatusCurrent)
	}

	// check decision
	if fetchedDP.Decision != nil && !toolbox.ArePointersEqual(fetchedDP.Decision, storedDP.Decision) {
		hasUpdate = true
		updateMessage += fmt.Sprintf("Decision updated to '%s'\n", *fetchedDP.Decision)
		if *&fetchedDP.DecisionBy != nil {
			updateMessage += fmt.Sprintf("Decision By '%s'\n", *fetchedDP.DecisionBy)
		}
	}

	return hasUpdate, updateMessage
}

// UpsertDevelopmentPermit - updates or inserts a development permet to a list of permits
func UpsertDevelopmentPermit(permits []DevelopmentPermit, thePermit DevelopmentPermit) []DevelopmentPermit {
	//Search the permits for the index of the permit to add. If found update, if not append
	foundIndex := slices.IndexFunc(permits, func(c DevelopmentPermit) bool { return c.PermitNum == thePermit.PermitNum })
	if foundIndex != -1 {
		permits[foundIndex] = thePermit
	} else {
		permits = append(permits, thePermit)
	}

	return permits
}
