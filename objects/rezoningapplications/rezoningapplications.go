package rezoningapplications

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/jeffadavidson/development-bot/interactions/calgaryopendata"
	"github.com/jeffadavidson/development-bot/interactions/rssfeed"
	"github.com/jeffadavidson/development-bot/objects/fileaction"
	"github.com/jeffadavidson/development-bot/utilities/fileio"
	"github.com/jeffadavidson/development-bot/utilities/toolbox"
	"golang.org/x/exp/slices"
)

// generateUniqueGUID creates a unique GUID based on the permit number and type
func generateUniqueGUID(permitNum, permitType string) string {
	// Create a deterministic GUID by hashing the permit number + type
	hash := sha256.Sum256([]byte(permitType + ":" + permitNum))
	return hex.EncodeToString(hash[:16]) // Use first 16 bytes for a shorter GUID
}

// updateStateHistory tracks state changes for rezoning applications
func updateStateHistory(fetchedApp *RezoningApplication, storedApp *RezoningApplication) {
	// Initialize state history from stored app if it exists
	if storedApp != nil {
		fetchedApp.StateHistory = storedApp.StateHistory
	}

	// Check if this is a new state we haven't seen before
	currentStatus := strings.ToLower(strings.TrimSpace(fetchedApp.StatusCurrent))
	timestamp := time.Now().Format(time.RFC3339)

	// If this is the first time we're seeing this application, add initial state
	if len(fetchedApp.StateHistory) == 0 {
		fetchedApp.StateHistory = append(fetchedApp.StateHistory, StateChange{
			Status:    currentStatus,
			Timestamp: timestamp,
		})
		return
	}

	// Check if the status has changed from the last recorded state
	lastState := fetchedApp.StateHistory[len(fetchedApp.StateHistory)-1]
	if strings.ToLower(strings.TrimSpace(lastState.Status)) != currentStatus {
		fetchedApp.StateHistory = append(fetchedApp.StateHistory, StateChange{
			Status:    currentStatus,
			Timestamp: timestamp,
		})
	}
}

// GetStateHistorySummary returns a human-readable summary of the application's lifecycle
func (ra RezoningApplication) GetStateHistorySummary() string {
	if len(ra.StateHistory) == 0 {
		return "No state history available"
	}
	
	summary := fmt.Sprintf("Application %s lifecycle:\n", ra.PermitNum)
	for i, state := range ra.StateHistory {
		timestamp, _ := time.Parse(time.RFC3339, state.Timestamp)
		summary += fmt.Sprintf("  %d. %s - %s\n", i+1, 
			strings.Title(state.Status), 
			timestamp.Format("Jan 2, 2006 3:04 PM"))
	}
	return summary
}

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
	RSSGuid                string         `json:"rss_guid"`
	StateHistory           []StateChange `json:"state_history"`
}

type StateChange struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
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

// generateRSSDescription creates a self-contained HTML description for RSS feeds
func (ra *RezoningApplication) generateRSSDescription() string {
	var html strings.Builder
	
	// Header with application number and status
	html.WriteString(fmt.Sprintf("<h3>🏛️ REZONING APPLICATION %s</h3>", ra.PermitNum))
	html.WriteString(fmt.Sprintf("<p><strong>Status:</strong> %s</p>", ra.StatusCurrent))
	
	// Address and location details with map links
	if ra.Address != nil {
		html.WriteString(fmt.Sprintf("<p>📍 <strong>Address:</strong> %s</p>", *ra.Address))
		
		// Map links right under address
		html.WriteString("<ul style='margin-top: 5px; margin-bottom: 15px;'>")
		googleMapsURL := fmt.Sprintf("https://maps.google.com/?q=%s", url.QueryEscape(fmt.Sprintf("%s, Calgary, Alberta", *ra.Address)))
		html.WriteString(fmt.Sprintf("<li>📍 <a href='%s' target='_blank'>View on Google Maps</a></li>", googleMapsURL))
		dmapURL := "https://developmentmap.calgary.ca/?find=" + ra.PermitNum
		html.WriteString(fmt.Sprintf("<li>📋 <a href='%s' target='_blank'>View on Calgary Development Map</a></li>", dmapURL))
		html.WriteString("</ul>")
	}
	
	// Project details
	if ra.Description != nil {
		html.WriteString(fmt.Sprintf("<p>🏗️ <strong>Project:</strong> %s</p>", *ra.Description))
	}
	
	// Land use change details
	html.WriteString("<div style='background-color: #fff3cd; padding: 10px; margin: 10px 0; border-left: 4px solid #ffc107;'>")
	html.WriteString("<h4>🏘️ LAND USE CHANGE:</h4>")
	html.WriteString("<ul>")
	
	if ra.FromLud != nil {
		html.WriteString(fmt.Sprintf("<li><strong>From:</strong> %s</li>", *ra.FromLud))
	}
	
	if ra.ProposedLud != nil {
		html.WriteString(fmt.Sprintf("<li><strong>To:</strong> %s</li>", *ra.ProposedLud))
	}
	html.WriteString("</ul></div>")
	
	// Applicant information
	if ra.Applicant != nil {
		html.WriteString(fmt.Sprintf("<p>👤 <strong>Applicant:</strong> %s</p>", *ra.Applicant))
	}
	
	// Timeline information
	html.WriteString("<h4>📅 TIMELINE:</h4><ul>")
	
	if ra.AppliedDate != nil {
		if parsedDate, err := time.Parse("2006-01-02T15:04:05.000", *ra.AppliedDate); err == nil {
			html.WriteString(fmt.Sprintf("<li>Applied: %s</li>", parsedDate.Format("January 2, 2006")))
		}
	}
	
	if ra.CompletedDate != nil && *ra.CompletedDate != "" {
		if parsedDate, err := time.Parse("2006-01-02T15:04:05.000", *ra.CompletedDate); err == nil {
			html.WriteString(fmt.Sprintf("<li>Completed: %s</li>", parsedDate.Format("January 2, 2006")))
		}
	}
	html.WriteString("</ul>")
	
	// Status information
	if ra.StatusCurrent != "" && strings.Contains(strings.ToLower(ra.StatusCurrent), "approv") {
		html.WriteString("<div style='background-color: #d4edda; padding: 10px; margin: 10px 0; border-left: 4px solid #28a745;'>")
		html.WriteString(fmt.Sprintf("<strong>✅ STATUS:</strong> %s", ra.StatusCurrent))
		html.WriteString("</div>")
	}
	

	
	return html.String()
}

// EvaluateRezoningApplications - Evaluates rezoning applications and generates RSS feed
func EvaluateRezoningApplications(rss *rssfeed.RSS) ([]fileaction.FileAction, error) {
	// Load rezoning applications
	fetchedPermits, storedPermits, err := loadRezoningApplications()
	if err != nil {
		return nil, fmt.Errorf("failed to load rezoning applications: %v", err)
	}
	fileActions := getRezoningApplicationActions(fetchedPermits, storedPermits)

	// Process actions for Rezoning Applications
	for _, val := range fileActions {
		if val.Action == "CREATE" {
			// Add new RSS item
			ra := findRezoningApplicationByID(fetchedPermits, val.PermitNum)
			if ra != nil {
				// Use the most recent timestamp from application data
				pubDate := ra.getMostRecentTimestamp()

				// Generate consistent title without status
				title := fmt.Sprintf("🏛️ Rezoning Application: %s", val.PermitNum)
				if ra.Address != nil {
					title = fmt.Sprintf("🏛️ Rezoning Application: %s - %s", val.PermitNum, *ra.Address)
				}

				link := fmt.Sprintf("https://developmentmap.calgary.ca/?find=%s", val.PermitNum)
				
				// Enhanced RSS metadata
				category := "Land Use Rezoning"
				author := "Unknown"
				if ra.Applicant != nil {
					author = *ra.Applicant
				}
				source := "City of Calgary Open Data"
				comments := fmt.Sprintf("https://developmentmap.calgary.ca/?find=%s#comments", val.PermitNum)
				
				// Use full content in both description and content:encoded for maximum compatibility
				fullContent := ra.generateRSSDescription()
				
				// Only update RSS and print messages if actual changes were made
				wasUpdated := rss.UpdateItem(title, fullContent, link, ra.RSSGuid, pubDate, category, author, source, comments, fullContent)
				if wasUpdated {
					fmt.Printf("Rezoning Application %s:\n\tCreating RSS feed entry...\n", val.PermitNum)
					fmt.Printf("\tCreated RSS feed entry!\n")
				}
			}
		}

		if val.Action == "UPDATE" || val.Action == "CLOSE" {
			// Update existing RSS item
			ra := findRezoningApplicationByID(fetchedPermits, val.PermitNum)
			if ra != nil {
				// Use the most recent timestamp from application data
				pubDate := ra.getMostRecentTimestamp()

				// Generate consistent title without status
				title := fmt.Sprintf("🏛️ Rezoning Application: %s", val.PermitNum)
				if ra.Address != nil {
					title = fmt.Sprintf("🏛️ Rezoning Application: %s - %s", val.PermitNum, *ra.Address)
				}

				link := fmt.Sprintf("https://developmentmap.calgary.ca/?find=%s", val.PermitNum)
				
				// Enhanced RSS metadata
				category := "Land Use Rezoning"
				author := "Unknown"
				if ra.Applicant != nil {
					author = *ra.Applicant
				}
				source := "City of Calgary Open Data"
				comments := fmt.Sprintf("https://developmentmap.calgary.ca/?find=%s#comments", val.PermitNum)
				
				// Use full content in both description and content:encoded for maximum compatibility
				fullContent := ra.generateRSSDescription()
				
				// Only update RSS and print messages if actual changes were made
				wasUpdated := rss.UpdateItem(title, fullContent, link, ra.RSSGuid, pubDate, category, author, source, comments, fullContent)
				if wasUpdated {
					fmt.Printf("Rezoning Application %s:\n\tUpdating RSS feed entry...\n", val.PermitNum)
					fmt.Printf("\tUpdated RSS feed entry!\n")
				}
			}
		}
	}

	// Save Rezoning Applications (save the fetched data so we can compare next time)
	saveRezoningApplications(fetchedPermits)

	return fileActions, nil
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

	// Ensure all fetched permits have GUIDs (generate if new, preserve if existing)
	for i := range fetchedRezoningApplications {
		storedPermit := findRezoningApplicationByID(storedPermits, fetchedRezoningApplications[i].PermitNum)
		if storedPermit != nil && storedPermit.RSSGuid != "" {
			// Use existing GUID from stored data
			fetchedRezoningApplications[i].RSSGuid = storedPermit.RSSGuid
		} else {
			// Generate new GUID for new permit
			fetchedRezoningApplications[i].RSSGuid = generateUniqueGUID(fetchedRezoningApplications[i].PermitNum, "rezoning-application")
		}

		// Update state history if status changed
		updateStateHistory(&fetchedRezoningApplications[i], storedPermit)
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

// getMostRecentTimestamp finds the most recent timestamp from a rezoning application's data
func (ra *RezoningApplication) getMostRecentTimestamp() time.Time {
	var mostRecent time.Time
	
	// Check applied date
	if ra.AppliedDate != nil {
		if appliedDate, err := time.Parse("2006-01-02T15:04:05.000", *ra.AppliedDate); err == nil {
			if appliedDate.After(mostRecent) {
				mostRecent = appliedDate
			}
		}
	}
	
	// Check completed date
	if ra.CompletedDate != nil {
		if completedDate, err := time.Parse("2006-01-02T15:04:05.000", *ra.CompletedDate); err == nil {
			if completedDate.After(mostRecent) {
				mostRecent = completedDate
			}
		}
	}
	
	// Check state history for most recent status change
	for _, state := range ra.StateHistory {
		if stateTime, err := time.Parse(time.RFC3339, state.Timestamp); err == nil {
			if stateTime.After(mostRecent) {
				mostRecent = stateTime
			}
		}
	}
	
	// If no valid timestamp found, use current time
	if mostRecent.IsZero() {
		mostRecent = time.Now()
	}
	
	return mostRecent
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

	// Check for close statuses
	close_statuses := [3]string{"Approved", "Cancelled", "Refused"}
	currentlyInCloseStatus := toolbox.SliceContains([]string(close_statuses[:]), fetchedRA.StatusCurrent)
	previouslyInCloseStatus := toolbox.SliceContains([]string(close_statuses[:]), storedRA.StatusCurrent)
	
	// Only close if currently in close status AND wasn't previously in close status
	if currentlyInCloseStatus && !previouslyInCloseStatus {
		toClose = true
		closeMessage = fmt.Sprintf("Closing file as it changed to status '%s'", fetchedRA.StatusCurrent)
	} else if currentlyInCloseStatus && previouslyInCloseStatus {
		// No action needed if already in closed status
	}

	return toClose, closeMessage
}
