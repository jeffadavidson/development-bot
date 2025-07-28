package rezoningapplications

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/jeffadavidson/development-bot/interactions/calgaryopendata"
	"github.com/jeffadavidson/development-bot/interactions/rssfeed"
	"github.com/jeffadavidson/development-bot/objects/fileaction"
	"github.com/jeffadavidson/development-bot/utilities/config"
	"github.com/jeffadavidson/development-bot/utilities/fileio"
	"github.com/jeffadavidson/development-bot/utilities/toolbox"
	"golang.org/x/exp/slices"
)

type RezoningApplication struct {
	PermitType             string     `json:"permittype"`
	PermitNum              string     `json:"permitnum"`
	Description            *string    `json:"description"`
	StatusCurrent          string     `json:"statuscurrent"`
	AppliedDate            *string    `json:"applieddate"`
	CompletedDate          *string    `json:"completeddate"`
	Applicant              *string    `json:"applicant"`
	FromLud                *string    `json:"fromlud"`
	ProposedLud            *string    `json:"proposedlud"`
	Address                *string    `json:"address"`
	LocationAddresses      *string    `json:"locationaddresses"`
	LocationCount          *string    `json:"locationcount"`
	Latitude               *string    `json:"latitude"`
	Longitude              *string    `json:"longitude"`
	Multipoint             Multipoint `json:"multipoint"`
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

// EvaluateRezoningApplications - Evaluates rezoning applications and generates RSS feed
func EvaluateRezoningApplications() error {
	// Load rezoning applications
	fetchedPermits, storedPermits, err := loadRezoningApplications()
	if err != nil {
		return fmt.Errorf("failed to load rezoning applications: %v", err)
	}
	fileActions := getRezoningApplicationActions(fetchedPermits, storedPermits)

	// Load or create RSS feed
	rss, err := rssfeed.GetOrCreateRSSFeed(
		"./data/rezoning-applications.xml",
		"Killarney Rezoning Applications",
		"Land use rezoning application updates for the Killarney neighborhood in Calgary",
		"https://calgary.ca/rezoning-applications",
	)
	if err != nil {
		return fmt.Errorf("failed to load RSS feed: %v", err)
	}

	// Process actions for Rezoning Applications
	for _, val := range fileActions {
		if val.Action == "CREATE" {
			fmt.Printf("Rezoning Application %s:\n\tAdding to RSS feed...\n", val.PermitNum)
			
			// Add new RSS item
			ra := findRezoningApplicationByID(fetchedPermits, val.PermitNum)
			if ra != nil {
				pubDate := time.Now()
				if ra.AppliedDate != nil {
					if parsedDate, parseErr := time.Parse("2006-01-02T15:04:05.000", *ra.AppliedDate); parseErr == nil {
						pubDate = parsedDate
					}
				}
				
				title := fmt.Sprintf("New Rezoning Application: %s", val.PermitNum)
				if ra.Address != nil {
					title = fmt.Sprintf("New Rezoning Application: %s - %s", val.PermitNum, *ra.Address)
				}
				
				rss.AddItem(title, val.Message, "", val.PermitNum, pubDate)
				fmt.Printf("\tAdded to RSS feed!\n")
			}
		}

		if val.Action == "UPDATE" || val.Action == "CLOSE" {
			fmt.Printf("Rezoning Application %s:\n\tUpdating RSS feed...\n", val.PermitNum)
			
			// Update existing RSS item
			ra := findRezoningApplicationByID(fetchedPermits, val.PermitNum)
			if ra != nil {
				pubDate := time.Now()
				
				title := fmt.Sprintf("Rezoning Application Update: %s", val.PermitNum)
				if ra.Address != nil {
					title = fmt.Sprintf("Rezoning Application Update: %s - %s", val.PermitNum, *ra.Address)
				}
				
				if val.Action == "CLOSE" {
					title = fmt.Sprintf("Rezoning Application Closed: %s", val.PermitNum)
					if ra.Address != nil {
						title = fmt.Sprintf("Rezoning Application Closed: %s - %s", val.PermitNum, *ra.Address)
					}
				}
				
				rss.UpdateItem(title, val.Message, "", val.PermitNum, pubDate)
				fmt.Printf("\tUpdated in RSS feed!\n")
			}
		}
	}

	// Trim RSS feed to keep only recent items
	rss.TrimToMaxItems(100)

	// Save RSS feed
	if err := rssfeed.SaveRSSFeed(rss, "./data/rezoning-applications.xml"); err != nil {
		return fmt.Errorf("failed to save RSS feed: %v", err)
	}

	// Save Rezoning Applications
	saveRezoningApplications(storedPermits)

	return nil
}

// loadRezoningApplications - Loads existing rezoning applications and fetches new ones from Calgary Open Data
func loadRezoningApplications() ([]RezoningApplication, []RezoningApplication, error) {
	// Load existing rezoning applications
	storedPermitsBytes, loadErr := fileio.GetFileContents("./data/rezoning-applications.json")
	if loadErr != nil {
		return nil, nil, loadErr
	}
	storedPermits, parseErr := parseRezoningApplications(storedPermitsBytes)
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

	return fetchedRezoningApplications, storedPermits, nil
}

// saveRezoningApplications - saves rezoning applications to file
func saveRezoningApplications(applications []RezoningApplication) error {
	// Encode applications as JSON
	applicationsBytes, encodeErr := json.MarshalIndent(applications, "", "  ")
	if encodeErr != nil {
		return encodeErr
	}

	writeErr := fileio.WriteFileContents("./data/rezoning-applications.json", applicationsBytes)
	if writeErr != nil {
		return writeErr
	}

	return nil
}

// findRezoningApplicationByID - finds a rezoning application in a list of applications
func findRezoningApplicationByID(searchSlice []RezoningApplication, id string) *RezoningApplication {
	foundIndex := slices.IndexFunc(searchSlice, func(c RezoningApplication) bool { return c.PermitNum == id })
	if foundIndex == -1 {
		return nil
	}

	return &searchSlice[foundIndex]
}



// getRezoningApplicationActions - Compares fetched and stored rezoning applications and returns a list of actions to execute
func getRezoningApplicationActions(fetchedRezoningApplications []RezoningApplication, storedPermits []RezoningApplication) []fileaction.FileAction {
	var fileActions []fileaction.FileAction
	for _, fetchedRA := range fetchedRezoningApplications {
		storedRAPtr := findRezoningApplicationByID(storedPermits, fetchedRA.PermitNum)

		if storedRAPtr == nil {
			// New application - create RSS entry
			fileActions = append(fileActions, fileaction.FileAction{PermitNum: fetchedRA.PermitNum, Action: "CREATE", Message: fetchedRA.CreateInformationMessage()})
		} else {
			storedRA := *storedRAPtr

			hasUpdate, updateMessage := getRezoningApplicationUpdates(fetchedRA, storedRA)
			toClose, closeMessage := isRezoningApplicationClosed(fetchedRA, storedRA)
			var message string

			if hasUpdate && !toClose {
				message += updateMessage
				fileActions = append(fileActions, fileaction.FileAction{PermitNum: fetchedRA.PermitNum, Action: "UPDATE", Message: message})
			}
			if hasUpdate && toClose {
				message += updateMessage
				message += "\n"
				message += closeMessage

				fileActions = append(fileActions, fileaction.FileAction{PermitNum: fetchedRA.PermitNum, Action: "CLOSE", Message: message})
			}
			if !hasUpdate && toClose {
				message += closeMessage
				fileActions = append(fileActions, fileaction.FileAction{PermitNum: fetchedRA.PermitNum, Action: "CLOSE", Message: closeMessage})
			}
		}
	}

	return fileActions
}

// upsertRezoningApplication - updates or inserts a rezoning application to a list of applications
func upsertRezoningApplication(applications []RezoningApplication, theApplication RezoningApplication) []RezoningApplication {
	// Search the applications for the index of the application to add. If found update, if not append
	foundIndex := slices.IndexFunc(applications, func(c RezoningApplication) bool { return c.PermitNum == theApplication.PermitNum })
	if foundIndex != -1 {
		applications[foundIndex] = theApplication
	} else {
		applications = append(applications, theApplication)
	}

	return applications
}

// parseRezoningApplications - Parses raw bytes from Calgary Open Data into an array of RezoningApplication objects
func parseRezoningApplications(raw []byte) ([]RezoningApplication, error) {
	var applications []RezoningApplication
	err := json.Unmarshal(raw, &applications)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rezoning application json. Error: %s", err.Error())
	}
	return applications, nil
}

// getRezoningApplicationUpdates - Checks if a rezoning application needs updates
func getRezoningApplicationUpdates(fetchedRA RezoningApplication, storedRA RezoningApplication) (bool, string) {
	hasUpdate := false
	updateMessage := ""

	// check status
	if fetchedRA.StatusCurrent != storedRA.StatusCurrent {
		hasUpdate = true
		updateMessage += fmt.Sprintf("Status updated from '%s' to '%s'\n", storedRA.StatusCurrent, fetchedRA.StatusCurrent)
	}

	return hasUpdate, updateMessage
}

// isRezoningApplicationClosed - Checks if a rezoning application is ready to be closed
func isRezoningApplicationClosed(fetchedRA RezoningApplication, storedRA RezoningApplication) (bool, string) {
	toClose := false
	closeMessage := ""

	// Check for close
	close_statuses := [3]string{"Approved", "Cancelled", "Refused"}
	if toolbox.SliceContains([]string(close_statuses[:]), fetchedRA.StatusCurrent) {
		toClose = true
		closeMessage = fmt.Sprintf("Closing file as it is in status '%s'", fetchedRA.StatusCurrent)
	}

	return toClose, closeMessage
}
