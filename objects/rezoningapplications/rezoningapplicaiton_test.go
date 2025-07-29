package rezoningapplications

import (
	"testing"
	"time"

	"github.com/jeffadavidson/development-bot/objects/fileaction"
	"github.com/stretchr/testify/assert"
)

func Test_FindRezoningApplication_ItemFound(t *testing.T) {
	dp1 := RezoningApplication{PermitNum: "123"}
	dp2 := RezoningApplication{PermitNum: "456"}
	dp3 := RezoningApplication{PermitNum: "789"}
	searchSlice := []RezoningApplication{dp1, dp2, dp3}

	result := findRezoningApplicationByID(searchSlice, "456")

	assert.NotNil(t, result)
	assert.Equal(t, &dp2, result)
}

func Test_FindRezoningApplication_ItemNotFound(t *testing.T) {
	dp1 := RezoningApplication{PermitNum: "123"}
	//	dp2 := RezoningApplication{PermitNum: "456"}
	dp3 := RezoningApplication{PermitNum: "789"}
	searchSlice := []RezoningApplication{dp1, dp3}

	result := findRezoningApplicationByID(searchSlice, "456")

	assert.Nil(t, result)
}

func Test_FindRezoningApplication_EmptySlice(t *testing.T) {
	searchSlice := []RezoningApplication{}

	result := findRezoningApplicationByID(searchSlice, "456")

	assert.Nil(t, result)
}

func Test_ParseRezoningApplication_Valid(t *testing.T) {
	fetchedRa := []byte(`
	[
		{
			"permittype": "Rezoning",
			"permitnum": "RZ2023-00001",
			"description": "Rezoning application for a new residential development",
			"statuscurrent": "In Progress",
			"applieddate": "2023-03-31",
			"completeddate": null,
			"applicant": "ABC Development Corp.",
			"fromlud": "R-C2",
			"proposedlud": "R-C1",
			"address": "123 Main St",
			"locationaddresses": "123 Main St",
			"locationcount": "1",
			"latitude": "51.123",
			"longitude": "-114.123",
			"multipoint": {
				"type": "MultiPoint",
				"coordinates": [
					[-114.123, 51.123]
				]
			},
			"github_discussion_id": null,
			"github_discussion_closed": false
		}
	]
	`)

	returned, err := parseRezoningApplications(fetchedRa)
	assert.NoError(t, err)
	assert.Len(t, returned, 1)
	assert.Equal(t, "RZ2023-00001", returned[0].PermitNum)
}

func Test_ParseRezoningApplication_MalformedJson(t *testing.T) {
	fetchedRa := []byte(`
	[
		{
			"permittype": "Rezoning",
			"permitnum": "RZ2023-00001",
			"description": "Rezoning application for a new residential development",
			"statuscurrent": "In Progress",
			"applieddate": "2023-03-31",
			"completeddate": null,
			"applicant": "ABC Development Corp.",
			"fromlud": "R-C2",
			"proposedlud": "R-C1",
			"address": "123 Main St",
			"locationaddresses": "123 Main St",
			"locationcount": "1",
			"latitude": "51.123",
			"longitude": "-114.123",
			"multipoint": {
				"type": "MultiPoint",
				"coordinates": [
					[-114.123, 51.123]
				]
			},
			"github_discussion_id": null,
			"github_discussion_closed": false
		
	]
	`)

	returned, err := parseRezoningApplications(fetchedRa)
	assert.NotEqual(t, nil, err)
	assert.ErrorContains(t, err, "failed to parse rezoning application json")
	assert.Equal(t, 0, len(returned))
}

func Test_RezoningApplication_NoUpdate(t *testing.T) {
	fetchedRaJson := []byte(`
	[
		{
			"permittype": "Rezoning",
			"permitnum": "RZ2023-00001",
			"description": "Rezoning application for a new residential development",
			"statuscurrent": "In Progress",
			"applieddate": "2023-03-31",
			"completeddate": null,
			"applicant": "ABC Development Corp.",
			"fromlud": "R-C2",
			"proposedlud": "R-C1",
			"address": "123 Main St",
			"locationaddresses": "123 Main St",
			"locationcount": "1",
			"latitude": "51.123",
			"longitude": "-114.123",
			"multipoint": {
				"type": "MultiPoint",
				"coordinates": [
					[-114.123, 51.123]
				]
			},
			"github_discussion_id": "12345",
			"github_discussion_closed": false
		}
	]
	`)
	storedRaJson := []byte(`
	[
		{
			"permittype": "Rezoning",
			"permitnum": "RZ2023-00001",
			"description": "Rezoning application for a new residential development",
			"statuscurrent": "In Progress",
			"applieddate": "2023-03-31",
			"completeddate": null,
			"applicant": "ABC Development Corp.",
			"fromlud": "R-C2",
			"proposedlud": "R-C1",
			"address": "123 Main St",
			"locationaddresses": "123 Main St",
			"locationcount": "1",
			"latitude": "51.123",
			"longitude": "-114.123",
			"multipoint": {
				"type": "MultiPoint",
				"coordinates": [
					[-114.123, 51.123]
				]
			},
			"github_discussion_id": "12345",
			"github_discussion_closed": false
		}
	]
	`)

	fetchedRa, errP := parseRezoningApplications(fetchedRaJson)
	storedRa, errE := parseRezoningApplications(storedRaJson)

	assert.NoError(t, errP)
	assert.NoError(t, errE)

	createdActions := getRezoningApplicationActions(fetchedRa, storedRa)
	assert.Equal(t, 0, len(createdActions))
}

func Test_RezoningApplication_StatusUpdate(t *testing.T) {
	fetchedRaJson := []byte(`
	[
		{
			"permittype": "Rezoning",
			"permitnum": "RZ2023-00001",
			"description": "Rezoning application for a new residential development",
			"statuscurrent": "In Progress",
			"applieddate": "2023-03-31",
			"completeddate": null,
			"applicant": "ABC Development Corp.",
			"fromlud": "R-C2",
			"proposedlud": "R-C1",
			"address": "123 Main St",
			"locationaddresses": "123 Main St",
			"locationcount": "1",
			"latitude": "51.123",
			"longitude": "-114.123",
			"multipoint": {
				"type": "MultiPoint",
				"coordinates": [
					[-114.123, 51.123]
				]
			},
			"github_discussion_id": "12345",
			"github_discussion_closed": false
		}
	]
	`)
	storedRaJson := []byte(`
	[
		{
			"permittype": "Rezoning",
			"permitnum": "RZ2023-00001",
			"description": "Rezoning application for a new residential development",
			"statuscurrent": "In Circulation",
			"applieddate": "2023-03-31",
			"completeddate": null,
			"applicant": "ABC Development Corp.",
			"fromlud": "R-C2",
			"proposedlud": "R-C1",
			"address": "123 Main St",
			"locationaddresses": "123 Main St",
			"locationcount": "1",
			"latitude": "51.123",
			"longitude": "-114.123",
			"multipoint": {
				"type": "MultiPoint",
				"coordinates": [
					[-114.123, 51.123]
				]
			},
			"github_discussion_id": "12345",
			"github_discussion_closed": false
		}
	]
	`)
	expectedActions := []fileaction.FileAction{
		{
			PermitNum: "RZ2023-00001",
			Action:    "UPDATE",
			Message:   "Status updated from 'In Circulation' to 'In Progress'\n",
		},
	}

	fetchedRa, errP := parseRezoningApplications(fetchedRaJson)
	storedRa, errE := parseRezoningApplications(storedRaJson)

	assert.NoError(t, errP)
	assert.NoError(t, errE)

	createdActions := getRezoningApplicationActions(fetchedRa, storedRa)
	assert.Equal(t, 1, len(createdActions))
	assert.Equal(t, expectedActions, createdActions)
}

func Test_RezoningApplication_StatusUpdateApproved(t *testing.T) {
	fetchedRaJson := []byte(`
	[
		{
			"permittype": "Rezoning",
			"permitnum": "RZ2023-00001",
			"description": "Rezoning application for a new residential development",
			"statuscurrent": "Approved",
			"applieddate": "2023-03-31",
			"completeddate": null,
			"applicant": "ABC Development Corp.",
			"fromlud": "R-C2",
			"proposedlud": "R-C1",
			"address": "123 Main St",
			"locationaddresses": "123 Main St",
			"locationcount": "1",
			"latitude": "51.123",
			"longitude": "-114.123",
			"multipoint": {
				"type": "MultiPoint",
				"coordinates": [
					[-114.123, 51.123]
				]
			},
			"github_discussion_id": "12345",
			"github_discussion_closed": false
		}
	]
	`)
	storedRaJson := []byte(`
	[
		{
			"permittype": "Rezoning",
			"permitnum": "RZ2023-00001",
			"description": "Rezoning application for a new residential development",
			"statuscurrent": "In Circulation",
			"applieddate": "2023-03-31",
			"completeddate": null,
			"applicant": "ABC Development Corp.",
			"fromlud": "R-C2",
			"proposedlud": "R-C1",
			"address": "123 Main St",
			"locationaddresses": "123 Main St",
			"locationcount": "1",
			"latitude": "51.123",
			"longitude": "-114.123",
			"multipoint": {
				"type": "MultiPoint",
				"coordinates": [
					[-114.123, 51.123]
				]
			},
			"github_discussion_id": "12345",
			"github_discussion_closed": false
		}
	]
	`)
	expectedActions := []fileaction.FileAction{
		{
			PermitNum: "RZ2023-00001",
			Action:    "CLOSE",
			Message:   "Status updated from 'In Circulation' to 'Approved'\n\nClosing file as it changed to status 'Approved'",
		},
	}

	fetchedRa, errP := parseRezoningApplications(fetchedRaJson)
	storedRa, errE := parseRezoningApplications(storedRaJson)

	assert.NoError(t, errP)
	assert.NoError(t, errE)

	createdActions := getRezoningApplicationActions(fetchedRa, storedRa)
	assert.Equal(t, 1, len(createdActions))
	assert.Equal(t, expectedActions, createdActions)
}

func Test_RezoningApplication_CreateRaDiscussion(t *testing.T) {
	fetchedRaJson := []byte(`
	[
		{
			"permittype": "Rezoning",
			"permitnum": "RZ2023-00001",
			"description": "Rezoning application for a new residential development",
			"statuscurrent": "Approved",
			"applieddate": "2023-03-31",
			"completeddate": null,
			"applicant": "ABC Development Corp.",
			"fromlud": "R-C2",
			"proposedlud": "R-C1",
			"address": "123 Main St",
			"locationaddresses": "123 Main St",
			"locationcount": "1",
			"latitude": "51.123",
			"longitude": "-114.123",
			"multipoint": {
				"type": "MultiPoint",
				"coordinates": [
					[-114.123, 51.123]
				]
			},
			"github_discussion_id": null,
			"github_discussion_closed": false
		}
	]
	`)
	storedRaJson := []byte(`[]`)
	expectedActions := []fileaction.FileAction{
		{
			PermitNum: "RZ2023-00001",
			Action:    "CREATE",
		},
	}

	fetchedRa, errP := parseRezoningApplications(fetchedRaJson)
	storedRa, errE := parseRezoningApplications(storedRaJson)

	assert.NoError(t, errP)
	assert.NoError(t, errE)

	createdActions := getRezoningApplicationActions(fetchedRa, storedRa)
	assert.Equal(t, 1, len(createdActions))
	assert.Equal(t, expectedActions[0].PermitNum, createdActions[0].PermitNum)
	assert.Equal(t, expectedActions[0].Action, createdActions[0].Action)
}

func TestRezoningApplicationActions_CreateVsUpdate(t *testing.T) {
	// Test CREATE action - new application not in stored data
	storedRA := []byte(`[]`)
	fetchedRA := []byte(`
[
	{
		"point": {
			"type": "Point",
			"coordinates": [-114.13992156438134, 51.02930819585162]
		},
		"permitnum": "LOC2025-12345",
		"address": "123 Test ST SW",
		"statuscurrent": "Under Review",
		"applieddate": "2025-01-01T00:00:00.000"
	}
]
`)

	storedApplications, errS := parseRezoningApplications(storedRA)
	fetchedApplications, errF := parseRezoningApplications(fetchedRA)

	assert.NoError(t, errS)
	assert.NoError(t, errF)

	actions := getRezoningApplicationActions(fetchedApplications, storedApplications)
	
	// Should generate CREATE action for new application
	assert.Equal(t, 1, len(actions))
	assert.Equal(t, "LOC2025-12345", actions[0].PermitNum)
	assert.Equal(t, "CREATE", actions[0].Action)
}

func TestRezoningApplicationActions_UpdateAction(t *testing.T) {
	// Test UPDATE action - existing application with status change
	storedRA := []byte(`
[
	{
		"point": {
			"type": "Point",
			"coordinates": [-114.13992156438134, 51.02930819585162]
		},
		"permitnum": "LOC2025-12345",
		"address": "123 Test ST SW",
		"statuscurrent": "Under Review",
		"applieddate": "2025-01-01T00:00:00.000"
	}
]
`)
	fetchedRA := []byte(`
[
	{
		"point": {
			"type": "Point",
			"coordinates": [-114.13992156438134, 51.02930819585162]
		},
		"permitnum": "LOC2025-12345",
		"address": "123 Test ST SW",
		"statuscurrent": "Approved",
		"applieddate": "2025-01-01T00:00:00.000",
		"decision": "Approved",
		"decisiondate": "2025-01-15T00:00:00.000"
	}
]
`)

	storedApplications, errS := parseRezoningApplications(storedRA)
	fetchedApplications, errF := parseRezoningApplications(fetchedRA)

	assert.NoError(t, errS)
	assert.NoError(t, errF)

	actions := getRezoningApplicationActions(fetchedApplications, storedApplications)
	
	// Should generate CLOSE action for status change to Approved
	assert.Equal(t, 1, len(actions))
	assert.Equal(t, "LOC2025-12345", actions[0].PermitNum)
	assert.Equal(t, "CLOSE", actions[0].Action)
	assert.Contains(t, actions[0].Message, "Status updated")
	assert.Contains(t, actions[0].Message, "Closing file")
}

func TestRezoningApplicationActions_NoActionWhenAlreadyClosed(t *testing.T) {
	// Test no action when application is already in closed status
	storedRA := []byte(`
[
	{
		"point": {
			"type": "Point",
			"coordinates": [-114.13992156438134, 51.02930819585162]
		},
		"permitnum": "LOC2025-12345",
		"address": "123 Test ST SW",
		"statuscurrent": "Approved",
		"applieddate": "2025-01-01T00:00:00.000",
		"decision": "Approved",
		"decisiondate": "2025-01-15T00:00:00.000"
	}
]
`)
	fetchedRA := []byte(`
[
	{
		"point": {
			"type": "Point",
			"coordinates": [-114.13992156438134, 51.02930819585162]
		},
		"permitnum": "LOC2025-12345",
		"address": "123 Test ST SW",
		"statuscurrent": "Approved",
		"applieddate": "2025-01-01T00:00:00.000",
		"decision": "Approved",
		"decisiondate": "2025-01-15T00:00:00.000"
	}
]
`)

	storedApplications, errS := parseRezoningApplications(storedRA)
	fetchedApplications, errF := parseRezoningApplications(fetchedRA)

	assert.NoError(t, errS)
	assert.NoError(t, errF)

	actions := getRezoningApplicationActions(fetchedApplications, storedApplications)
	
	// Should generate no actions since application is already closed and unchanged
	assert.Equal(t, 0, len(actions))
}

func TestGetMostRecentTimestamp_WithStateHistory(t *testing.T) {
	ra := RezoningApplication{
		PermitNum:     "LOC2025-12345",
		AppliedDate:   strPtr("2025-01-01T10:00:00.000"),
		CompletedDate: strPtr("2025-01-15T14:30:00.000"),
		StateHistory: []StateChange{
			{
				Status:    "under review",
				Timestamp: "2025-01-01T10:00:00-07:00",
			},
			{
				Status:    "approved",
				Timestamp: "2025-01-20T16:45:00-07:00", // This should be the most recent
			},
		},
	}

	mostRecent := ra.getMostRecentTimestamp()
	
	// Should use the most recent state history timestamp
	expected, _ := time.Parse(time.RFC3339, "2025-01-20T16:45:00-07:00")
	assert.Equal(t, expected.Unix(), mostRecent.Unix())
}

func TestGetMostRecentTimestamp_WithoutStateHistory(t *testing.T) {
	ra := RezoningApplication{
		PermitNum:     "LOC2025-12345",
		AppliedDate:   strPtr("2025-01-01T10:00:00.000"),
		CompletedDate: strPtr("2025-01-20T16:45:00.000"), // This should be the most recent
		StateHistory:  []StateChange{}, // Empty state history
	}

	mostRecent := ra.getMostRecentTimestamp()
	
	// Should use the completed date as the most recent
	expected, _ := time.Parse("2006-01-02T15:04:05.000", "2025-01-20T16:45:00.000")
	assert.Equal(t, expected.Unix(), mostRecent.Unix())
}

func TestGetMostRecentTimestamp_NoTimestamps(t *testing.T) {
	ra := RezoningApplication{
		PermitNum:    "LOC2025-12345",
		StateHistory: []StateChange{},
	}

	mostRecent := ra.getMostRecentTimestamp()
	
	// Should use current time when no timestamps are available
	now := time.Now()
	// Allow for a small time difference since the function calls time.Now()
	assert.WithinDuration(t, now, mostRecent, time.Second)
}

// Helper function for string pointers
func strPtr(s string) *string {
	return &s
}
