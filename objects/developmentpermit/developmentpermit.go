package developmentpermit

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

// updateStateHistory tracks state changes for development permits
func updateStateHistory(fetchedPermit *DevelopmentPermit, storedPermit *DevelopmentPermit) {
	// Initialize state history from stored permit if it exists
	if storedPermit != nil {
		fetchedPermit.StateHistory = storedPermit.StateHistory
	}

	// Check if this is a new state we haven't seen before
	currentStatus := strings.ToLower(strings.TrimSpace(fetchedPermit.StatusCurrent))
	timestamp := time.Now().Format(time.RFC3339)

	// If this is the first time we're seeing this permit, add initial state
	if len(fetchedPermit.StateHistory) == 0 {
		decision := ""
		if fetchedPermit.Decision != nil {
			decision = *fetchedPermit.Decision
		}
		fetchedPermit.StateHistory = append(fetchedPermit.StateHistory, StateChange{
			Status:    currentStatus,
			Timestamp: timestamp,
			Decision:  decision,
		})
		return
	}

	// Check if the status has changed from the last recorded state
	lastState := fetchedPermit.StateHistory[len(fetchedPermit.StateHistory)-1]
	if strings.ToLower(strings.TrimSpace(lastState.Status)) != currentStatus {
		decision := ""
		if fetchedPermit.Decision != nil {
			decision = *fetchedPermit.Decision
		}
		fetchedPermit.StateHistory = append(fetchedPermit.StateHistory, StateChange{
			Status:    currentStatus,
			Timestamp: timestamp,
			Decision:  decision,
		})
	}
}

// GetStateHistorySummary returns a human-readable summary of the permit's lifecycle
func (dp DevelopmentPermit) GetStateHistorySummary() string {
	if len(dp.StateHistory) == 0 {
		return "No state history available"
	}
	
	summary := fmt.Sprintf("Permit %s lifecycle:\n", dp.PermitNum)
	for i, state := range dp.StateHistory {
		timestamp, _ := time.Parse(time.RFC3339, state.Timestamp)
		summary += fmt.Sprintf("  %d. %s - %s", i+1, 
			strings.Title(state.Status), 
			timestamp.Format("Jan 2, 2006 3:04 PM"))
		if state.Decision != "" {
			summary += fmt.Sprintf(" (Decision: %s)", state.Decision)
		}
		summary += "\n"
	}
	return summary
}

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
	RSSGuid                string         `json:"rss_guid"`
	StateHistory           []StateChange `json:"state_history"`
}

type StateChange struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Decision  string `json:"decision,omitempty"`
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

// generateRSSDescription creates a self-contained HTML description for RSS feeds
func (dp *DevelopmentPermit) generateRSSDescription() string {
	var html strings.Builder
	
	// Header with permit number and status
	html.WriteString(fmt.Sprintf("<h3>🏗️ DEVELOPMENT PERMIT %s</h3>", dp.PermitNum))
	html.WriteString(fmt.Sprintf("<p><strong>Status:</strong> %s</p>", dp.StatusCurrent))
	
	// Address and location details
	if dp.Address != nil {
		html.WriteString(fmt.Sprintf("<p>📍 <strong>Address:</strong> %s</p>", *dp.Address))
	}
	
	if dp.CommunityName != nil {
		html.WriteString(fmt.Sprintf("<p>🏘️ <strong>Community:</strong> %s", *dp.CommunityName))
		if dp.Ward != nil {
			html.WriteString(fmt.Sprintf(" (Ward %s)", *dp.Ward))
		}
		html.WriteString("</p>")
	}
	
	// Project details
	if dp.Description != nil {
		html.WriteString(fmt.Sprintf("<p>🏗️ <strong>Project:</strong> %s</p>", *dp.Description))
	}
	
	if dp.LandUseDistrictDesc != nil {
		html.WriteString(fmt.Sprintf("<p>🏘️ <strong>Land Use:</strong> %s</p>", *dp.LandUseDistrictDesc))
	}
	
	if dp.PermittedDiscretion != nil {
		html.WriteString(fmt.Sprintf("<p>📋 <strong>Application Type:</strong> %s</p>", *dp.PermittedDiscretion))
	}
	
	// Applicant information
	if dp.Applicant != nil {
		html.WriteString(fmt.Sprintf("<p>👤 <strong>Applicant:</strong> %s</p>", *dp.Applicant))
	}
	
	// Timeline information
	html.WriteString("<h4>📅 TIMELINE:</h4><ul>")
	
	if dp.AppliedDate != nil {
		if parsedDate, err := time.Parse("2006-01-02T15:04:05.000", *dp.AppliedDate); err == nil {
			html.WriteString(fmt.Sprintf("<li>Applied: %s</li>", parsedDate.Format("January 2, 2006")))
		}
	}
	
	if dp.DecisionDate != nil && *dp.DecisionDate != "" {
		if parsedDate, err := time.Parse("2006-01-02T15:04:05.000", *dp.DecisionDate); err == nil {
			html.WriteString(fmt.Sprintf("<li>Decision: %s</li>", parsedDate.Format("January 2, 2006")))
		}
	}
	
	if dp.ReleaseDate != nil && *dp.ReleaseDate != "" {
		if parsedDate, err := time.Parse("2006-01-02T15:04:05.000", *dp.ReleaseDate); err == nil {
			html.WriteString(fmt.Sprintf("<li>Released: %s</li>", parsedDate.Format("January 2, 2006")))
		}
	}
	
	if dp.MustCommenceDate != nil && *dp.MustCommenceDate != "" {
		if parsedDate, err := time.Parse("2006-01-02T15:04:05.000", *dp.MustCommenceDate); err == nil {
			html.WriteString(fmt.Sprintf("<li>Must Commence By: %s</li>", parsedDate.Format("January 2, 2006")))
		}
	}
	html.WriteString("</ul>")
	
	// Decision information
	if dp.Decision != nil && *dp.Decision != "" {
		html.WriteString("<div style='background-color: #d4edda; padding: 10px; margin: 10px 0; border-left: 4px solid #28a745;'>")
		html.WriteString(fmt.Sprintf("<strong>✅ DECISION:</strong> %s", *dp.Decision))
		if dp.DecisionBy != nil && *dp.DecisionBy != "" {
			html.WriteString(fmt.Sprintf(" (by %s)", *dp.DecisionBy))
		}
		html.WriteString("</div>")
	}
	
	// Location coordinates (if available)
	if dp.Latitude != nil && dp.Longitude != nil {
		html.WriteString(fmt.Sprintf("<p>📍 <strong>Coordinates:</strong> %s, %s</p>", *dp.Latitude, *dp.Longitude))
	}
	
	// Clickable links section
	html.WriteString("<hr/>")
	html.WriteString("<h4>🗺️ MAPS & DETAILS:</h4>")
	html.WriteString("<ul>")
	
	// Google Maps link using address
	if dp.Address != nil {
		googleMapsURL := fmt.Sprintf("https://maps.google.com/?q=%s", url.QueryEscape(fmt.Sprintf("%s, Calgary, Alberta", *dp.Address)))
		html.WriteString(fmt.Sprintf("<li>📍 <a href='%s' target='_blank'>View on Google Maps</a></li>", googleMapsURL))
	}
	
	// Development Map link
	dmapURL := "https://developmentmap.calgary.ca/?find=" + dp.PermitNum
	html.WriteString(fmt.Sprintf("<li>📋 <a href='%s' target='_blank'>View on Calgary Development Map</a></li>", dmapURL))
	
	html.WriteString("</ul>")
	
	return html.String()
}

func EvaluateDevelopmentPermits(rss *rssfeed.RSS) ([]fileaction.FileAction, error) {
	fetchedDevelopmentPermits, storedDevelopmentPermits, err := loadDevelopmentPermits()
	if err != nil {
		return nil, err
	}
	fileActions := getDevelopmentPermitActions(fetchedDevelopmentPermits, storedDevelopmentPermits)

	// Process actions for Development Permits
	for _, val := range fileActions {
		if val.Action == "CREATE" {
			fmt.Printf("Development Permit %s:\n\tUpdating RSS feed...\n", val.PermitNum)
			
						// Add new RSS item
			dp := findDevelopmentPermitByPermitNum(fetchedDevelopmentPermits, val.PermitNum)
			if dp != nil {
				pubDate := time.Now()
				if dp.AppliedDate != nil {
					if parsedDate, parseErr := time.Parse("2006-01-02T15:04:05.000", *dp.AppliedDate); parseErr == nil {
						pubDate = parsedDate
					}
				}

				// Generate status-based title
				status := strings.Title(strings.ToLower(dp.StatusCurrent))
				title := fmt.Sprintf("🏗️ Development Permit (%s): %s", status, val.PermitNum)
				if dp.Address != nil {
					title = fmt.Sprintf("🏗️ Development Permit (%s): %s - %s", status, val.PermitNum, *dp.Address)
				}

				link := fmt.Sprintf("https://developmentmap.calgary.ca/?find=%s", val.PermitNum)
				
				// Enhanced RSS metadata
				category := "Development Permit"
				author := "Unknown"
				if dp.Applicant != nil {
					author = *dp.Applicant
				}
				source := "City of Calgary Open Data"
				comments := fmt.Sprintf("https://developmentmap.calgary.ca/?find=%s#comments", val.PermitNum)
				
				// Use full content and create summary
				fullContent := dp.generateRSSDescription()
				summary := fmt.Sprintf("Development permit for %s - Status: %s", *dp.Address, dp.StatusCurrent)
				
				rss.UpdateItem(title, summary, link, dp.RSSGuid, pubDate, category, author, source, comments, fullContent)
				fmt.Printf("\tUpdated RSS feed!\n")
			}
		}

		if val.Action == "UPDATE" || val.Action == "CLOSE" {
			fmt.Printf("Development Permit %s:\n\tUpdating RSS feed...\n", val.PermitNum)
			
						// Update existing RSS item
			dp := findDevelopmentPermitByPermitNum(fetchedDevelopmentPermits, val.PermitNum)
			if dp != nil {
				// Keep original publication date for updates
				pubDate := time.Now()
				if dp.AppliedDate != nil {
					if parsedDate, parseErr := time.Parse("2006-01-02T15:04:05.000", *dp.AppliedDate); parseErr == nil {
						pubDate = parsedDate
					}
				}

				// Generate status-based title (always show current status)
				status := strings.Title(strings.ToLower(dp.StatusCurrent))
				title := fmt.Sprintf("🏗️ Development Permit (%s): %s", status, val.PermitNum)
				if dp.Address != nil {
					title = fmt.Sprintf("🏗️ Development Permit (%s): %s - %s", status, val.PermitNum, *dp.Address)
				}

				link := fmt.Sprintf("https://developmentmap.calgary.ca/?find=%s", val.PermitNum)
				
				// Enhanced RSS metadata
				category := "Development Permit"
				author := "Unknown"
				if dp.Applicant != nil {
					author = *dp.Applicant
				}
				source := "City of Calgary Open Data"
				comments := fmt.Sprintf("https://developmentmap.calgary.ca/?find=%s#comments", val.PermitNum)
				
				// Use full content and create summary
				fullContent := dp.generateRSSDescription()
				summary := fmt.Sprintf("Development permit for %s - Status: %s", *dp.Address, dp.StatusCurrent)
				
				rss.UpdateItem(title, summary, link, dp.RSSGuid, pubDate, category, author, source, comments, fullContent)
				fmt.Printf("\tUpdated in RSS feed!\n")
			}
		}
	}

	// Save Development Permits (save the fetched data so we can compare next time)
	saveDevelopmentPermits(fetchedDevelopmentPermits)

	return fileActions, nil
}

// loadDevelopmentPermits - Gets fetched development permits, gets stored development permits
func loadDevelopmentPermits() ([]DevelopmentPermit, []DevelopmentPermit, error) {
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

	// Ensure all fetched permits have GUIDs (generate if new, preserve if existing)
	for i := range fetchedDevelopmentPermits {
		storedPermit := findDevelopmentPermitByPermitNum(storedDevelopmentPermits, fetchedDevelopmentPermits[i].PermitNum)
		if storedPermit != nil && storedPermit.RSSGuid != "" {
			// Use existing GUID from stored data
			fetchedDevelopmentPermits[i].RSSGuid = storedPermit.RSSGuid
		} else {
			// Generate new GUID for new permit
			fetchedDevelopmentPermits[i].RSSGuid = generateUniqueGUID(fetchedDevelopmentPermits[i].PermitNum, "development-permit")
		}

		// Update state history if status changed
		updateStateHistory(&fetchedDevelopmentPermits[i], storedPermit)
	}

	return fetchedDevelopmentPermits, storedDevelopmentPermits, nil
}

func saveDevelopmentPermits(permits []DevelopmentPermit) error {
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

// findDevelopmentPermitByPermitNum - finds a development permit in a list of permits
func findDevelopmentPermitByPermitNum(searchSlice []DevelopmentPermit, permitNum string) *DevelopmentPermit {
	foundIndex := slices.IndexFunc(searchSlice, func(c DevelopmentPermit) bool { return c.PermitNum == permitNum })
	if foundIndex == -1 {
		return nil
	}

	return &searchSlice[foundIndex]
}



// getDevelopmentPermitActions - For a list of fetched and stored development permits compares permits and gets a list of actions to execute
func getDevelopmentPermitActions(fetchedDevelopmentPermits []DevelopmentPermit, storedDevelopmentPermits []DevelopmentPermit) []fileaction.FileAction {
	var fileActions []fileaction.FileAction
	for _, fetchedDP := range fetchedDevelopmentPermits {
		storedDpPtr := findDevelopmentPermitByPermitNum(storedDevelopmentPermits, fetchedDP.PermitNum)
		if storedDpPtr == nil {
			// New permit - create RSS entry
			fileActions = append(fileActions, fileaction.FileAction{PermitNum: fetchedDP.PermitNum, Action: "CREATE", Message: fetchedDP.CreateInformationMessage()})
		} else {
			storedDP := *storedDpPtr

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

// upsertDevelopmentPermit - updates or inserts a development permet to a list of permits
func upsertDevelopmentPermit(permits []DevelopmentPermit, thePermit DevelopmentPermit) []DevelopmentPermit {
	//Search the permits for the index of the permit to add. If found update, if not append
	foundIndex := slices.IndexFunc(permits, func(c DevelopmentPermit) bool { return c.PermitNum == thePermit.PermitNum })
	if foundIndex != -1 {
		permits[foundIndex] = thePermit
	} else {
		permits = append(permits, thePermit)
	}

	return permits
}

// parseDevelopmentPermits - parses a json byte array into objects
func parseDevelopmentPermits(developmentPermitByte []byte) ([]DevelopmentPermit, error) {
	var developmentPermits []DevelopmentPermit
	err := json.Unmarshal(developmentPermitByte, &developmentPermits)
	if err != nil {
		return developmentPermits, fmt.Errorf("failed to parse development permit json. Error: %s", err.Error())
	}

	return developmentPermits, nil
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
