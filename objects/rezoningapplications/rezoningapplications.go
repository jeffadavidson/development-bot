package rezoningapplications

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/jeffadavidson/development-bot/interactions/calgaryopendata"
	"github.com/jeffadavidson/development-bot/utilities/fileio"
)

type RezoningApplication struct {
	PermitType        string     `json:"permittype"`
	PermitNum         string     `json:"permitnum"`
	Description       *string    `json:"description"`
	StatusCurrent     string     `json:"statuscurrent"`
	AppliedDate       *string    `json:"applieddate"`
	CompletedDate     *string    `json:"completeddate"`
	Applicant         *string    `json:"applicant"`
	FromLud           *string    `json:"fromlud"`
	ProposedLud       *string    `json:"proposedlud"`
	Address           *string    `json:"address"`
	LocationAddresses *string    `json:"locationaddresses"`
	LocationCount     *string    `json:"locationcount"`
	Latitude          *string    `json:"latitude"`
	Longitude         *string    `json:"longitude"`
	Multipoint        Multipoint `json:"multipoint"`
}

type Multipoint struct {
	Type        **string    `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

// CreateInformationMessage - Builds an information message from the development permit
func (ra RezoningApplication) CreateInformationMessage() string {
	var lineSeparator string = "\n"
	var message string = ""

	// Markdown
	message += "## About\n\n"
	message += fmt.Sprintf("**Permit Number:** %v", ra.PermitNum)
	if ra.AppliedDate != nil {
		appliedDate, err := time.Parse("2006-01-02T15:04:05.000", *ra.AppliedDate)
		if err != nil {
			message += fmt.Sprintf("%v**Date Applied:** %v", lineSeparator, *ra.AppliedDate)
		} else {
			dateStr := appliedDate.Format("2006-01-02")
			message += fmt.Sprintf("%v**Date Applied:** %v", lineSeparator, dateStr)
		}
	}
	if ra.Address != nil {
		message += fmt.Sprintf("%v**Address:** %v", lineSeparator, *ra.Address)
	}
	if ra.Applicant != nil {
		message += fmt.Sprintf("%v**Applicant:** %v", lineSeparator, *ra.Applicant)
	}
	if ra.Description != nil {
		message += fmt.Sprintf("%v**Description:** %v", lineSeparator, *ra.Description)
	}
	if ra.ProposedLud != nil {
		message += fmt.Sprintf("%v**Proposed Land Use District:** %v", lineSeparator, *ra.ProposedLud)
	}
	if ra.FromLud != nil {
		message += fmt.Sprintf("%v**Current Land Use District:** %v", lineSeparator, *ra.FromLud)
	}
	if ra.StatusCurrent != "" {
		message += fmt.Sprintf("%v**Permit Status:** %v", lineSeparator, ra.StatusCurrent)
	}
	if ra.Multipoint.Coordinates != nil {
		message += fmt.Sprintf("%v**Coordinates:** %v", lineSeparator, ra.Multipoint.Coordinates)
	}

	// Add google maps link
	message += "\n## Links\n\n"
	message += lineSeparator + lineSeparator
	message += fmt.Sprintf("%v[Development Map](https://dmap.calgary.ca/?find=%v)", lineSeparator, ra.PermitNum)
	if ra.Address != nil {
		message += fmt.Sprintf("%v [Google Maps](https://maps.google.com/?q=%v)", lineSeparator, url.QueryEscape(fmt.Sprintf("%v, Calgary, Alberta", *ra.Address)))
	}
	message += lineSeparator

	return message
}

// EvaluateDevelopmentPermits - Evaluates development permits for rezoning applications
func EvaluateDevelopmentPermits(repositoryID string, categoryID string) error {
	// Load development permits
	fetchedPermits, storedPermits, err := loadRezoningApplications()
	if err != nil {
		return fmt.Errorf("failed to load development permits: %v", err)
	}

	fmt.Println("------------------------------------------------------------------------------------------")
	fmt.Println(fetchedPermits)
	fmt.Println("------------------------------------------------------------------------------------------")
	os.Exit(1)
	fmt.Println(storedPermits)
	fmt.Println("------------------------------------------------------------------------------------------")

	// TODO: Implement evaluation logic

	return nil
}

// loadRezoningApplications - Loads existing rezoning applications and fetches new ones from Calgary Open Data
func loadRezoningApplications() ([]RezoningApplication, []RezoningApplication, error) {
	// Load existing rezoning applications
	storedRezoningApplicationsBytes, loadErr := fileio.GetFileContents("./data/rezoning-applications.json")
	if loadErr != nil {
		return nil, nil, loadErr
	}
	storedRezoningApplications, parseErr := parseRezoningApplications(storedRezoningApplicationsBytes)
	if parseErr != nil {
		return nil, nil, parseErr
	}

	// Get rezoning applications from Calgary Open Data
	fetchedRezoningApplicationsRaw, fetchErr := calgaryopendata.GetRezoningApplications()
	if fetchErr != nil {
		return nil, nil, fetchErr
	}
	fetchedRezoningApplications, parseErr2 := parseRezoningApplications(fetchedRezoningApplicationsRaw)
	if parseErr2 != nil {
		return nil, nil, parseErr
	}

	return fetchedRezoningApplications, storedRezoningApplications, nil
}

// parseRezoningApplications - Parses raw bytes from Calgary Open Data into an array of RezoningApplication objects
func parseRezoningApplications(raw []byte) ([]RezoningApplication, error) {
	var applications []RezoningApplication
	err := json.Unmarshal(raw, &applications)
	if err != nil {
		return nil, err
	}
	return applications, nil
}
