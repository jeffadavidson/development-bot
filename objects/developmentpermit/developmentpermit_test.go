package developmentpermit

import (
	"testing"
	"time"

	"github.com/jeffadavidson/development-bot/objects/fileaction"
	"github.com/stretchr/testify/assert"
)

func Test_FindDevelopmentPermit_ItemFound(t *testing.T) {
	dp1 := DevelopmentPermit{PermitNum: "123"}
	dp2 := DevelopmentPermit{PermitNum: "456"}
	dp3 := DevelopmentPermit{PermitNum: "789"}
	searchSlice := []DevelopmentPermit{dp1, dp2, dp3}

	result := findDevelopmentPermitByPermitNum(searchSlice, "456")

	assert.NotNil(t, result)
	assert.Equal(t, &dp2, result)
}

func Test_FindDevelopmentPermit_ItemNotFound(t *testing.T) {
	dp1 := DevelopmentPermit{PermitNum: "123"}
	//	dp2 := DevelopmentPermit{PermitNum: "456"}
	dp3 := DevelopmentPermit{PermitNum: "789"}
	searchSlice := []DevelopmentPermit{dp1, dp3}

	result := findDevelopmentPermitByPermitNum(searchSlice, "456")

	assert.Nil(t, result)
}

func Test_FindDevelopmentPermit_EmptySlice(t *testing.T) {
	searchSlice := []DevelopmentPermit{}

	result := findDevelopmentPermitByPermitNum(searchSlice, "456")

	assert.Nil(t, result)
}

func Test_ParseDevelopmentPermit_Valid(t *testing.T) {
	dpJson := []byte(`
	[
		{
			"point": {
				"type": "Point",
				"coordinates": [-114.13992156438134, 51.02930819585162]
			},
			"permitnum": "DP2022-00156",
			"address": "2819 36 ST SW",
			"applicant": "TRICOR DESIGN GROUP",
			"category": "Residential - New Single / Semi / Duplex",
			"description": "NEW: CONTEXTUAL SEMI-DETACHED DWELLING, ACCESSORY RESIDENTIAL BUILDING (GARAGE)",
			"proposedusecode": "C1020; C1368",
			"proposedusedescription": "ACCESSORY RESIDENTIAL BUILDING; CONTEXTUAL SEMI-DETACHED DWELLING",
			"permitteddiscretionary": "Permitted",
			"landusedistrict": "R-CG",
			"landusedistrictdescription": "Residential - Grade-Oriented Infill",
			"statuscurrent": "Released",
			"applieddate": "2022-01-10T00:00:00.000",
			"decisiondate": "2022-02-16T00:00:00.000",
			"releasedate": "2022-04-19T00:00:00.000",
			"mustcommencedate": "2024-02-16T00:00:00.000",
			"decision": "Approval",
			"decisionby": "Development Authority",
			"communitycode": "KIL",
			"communityname": "KILLARNEY/GLENGARRY",
			"ward": "8",
			"quadrant": "SW",
			"latitude": "51.029",
			"longitude": "-114.140",
			"locationcount": "2",
			"locationtypes": "Titled Parcel;Building",
			"locationaddresses": "2819 36 ST SW;2819 36 ST SW",
			"locationsgeojson": "{\"type\":\"MultiPoint\",\"coordinates\":[[-114.1399216,51.0293082],[-114.1400278,51.029353]]}",
			"locationswkt": "MULTIPOINT ((-114.13992156438134 51.02930819585162),(-114.14002777754924 51.029352999468266))"
		}
	]
	
`)

	permits, err := parseDevelopmentPermits(dpJson)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(permits))
}

func Test_ParseDevelopmentPermit_MalformedJson(t *testing.T) {
	dpJson := []byte(`
	[
		{
			"point": {
				"type": "Point",
				"coordinates": [-114.13992156438134, 51.02930819585162]
			},
			"permitnum": "DP2022-00156",
			"address": "2819 36 ST SW",
			"applicant": "TRICOR DESIGN GROUP",
			"category": "Residential - New Single / Semi / Duplex",
			"description": "NEW: CONTEXTUAL SEMI-DETACHED DWELLING, ACCESSORY RESIDENTIAL BUILDING (GARAGE)",
			"proposedusecode": "C1020; C1368",
			"proposedusedescription": "ACCESSORY RESIDENTIAL BUILDING; CONTEXTUAL SEMI-DETACHED DWELLING",
			"permitteddiscretionary": "Permitted",
			"landusedistrict": "R-CG",
			"landusedistrictdescription": "Residential - Grade-Oriented Infill",
			"statuscurrent": "Released",
			"applieddate": "2022-01-10T00:00:00.000",
			"decisiondate": "2022-02-16T00:00:00.000",
			"releasedate": "2022-04-19T00:00:00.000",
			"mustcommencedate": "2024-02-16T00:00:00.000",
			"decision": "Approval",
			"decisionby": "Development Authority",
			"communitycode": "KIL",
			"communityname": "KILLARNEY/GLENGARRY",
			"ward": "8",
			"quadrant": "SW",
			"latitude": "51.029",
			"longitude": "-114.140",
			"locationcount": "2",
			"locationtypes": "Titled Parcel;Building",
			"locationaddresses": "2819 36 ST SW;2819 36 ST SW",
			"locationsgeojson": "{\"type\":\"MultiPoint\",\"coordinates\":[[-114.1399216,51.0293082],[-114.1400278,51.029353]]}",
			"locationswkt": "MULTIPOINT ((-114.13992156438134 51.02930819585162),(-114.14002777754924 51.029352999468266))"
	]
	
`)

	permits, err := parseDevelopmentPermits(dpJson)
	assert.NotEqual(t, nil, err)
	assert.ErrorContains(t, err, "failed to parse development permit json")
	assert.Equal(t, 0, len(permits))
}

func TestContains_DevelopmentPermitAction_NoUpdate(t *testing.T) {
	storedDP := []byte(`
	[
		{
			"point": {
				"type": "Point",
				"coordinates": [-114.13992156438134, 51.02930819585162]
			},
			"permitnum": "DP2022-00156",
			"address": "2819 36 ST SW",
			"applicant": "TRICOR DESIGN GROUP",
			"category": "Residential - New Single / Semi / Duplex",
			"description": "NEW: CONTEXTUAL SEMI-DETACHED DWELLING, ACCESSORY RESIDENTIAL BUILDING (GARAGE)",
			"proposedusecode": "C1020; C1368",
			"proposedusedescription": "ACCESSORY RESIDENTIAL BUILDING; CONTEXTUAL SEMI-DETACHED DWELLING",
			"permitteddiscretionary": "Permitted",
			"landusedistrict": "R-CG",
			"landusedistrictdescription": "Residential - Grade-Oriented Infill",
			"statuscurrent": "In Circulation",
			"applieddate": "2022-01-10T00:00:00.000",
			"decisiondate": "2022-02-16T00:00:00.000",
			"releasedate": "2022-04-19T00:00:00.000",
			"mustcommencedate": "2024-02-16T00:00:00.000",
			"decision": "Approval",
			"decisionby": "Development Authority",
			"communitycode": "KIL",
			"communityname": "KILLARNEY/GLENGARRY",
			"ward": "8",
			"quadrant": "SW",
			"latitude": "51.029",
			"longitude": "-114.140",
			"locationcount": "2",
			"locationtypes": "Titled Parcel;Building",
			"locationaddresses": "2819 36 ST SW;2819 36 ST SW",
			"locationsgeojson": "{\"type\":\"MultiPoint\",\"coordinates\":[[-114.1399216,51.0293082],[-114.1400278,51.029353]]}",
			"locationswkt": "MULTIPOINT ((-114.13992156438134 51.02930819585162),(-114.14002777754924 51.029352999468266))",
			"github_discussion_id": "D_kwDOJSUT984ATT9k",
			"github_discussion_closed": false
		}
	]
`)
	fetchedDP := []byte(`
[
	{
		"point": {
			"type": "Point",
			"coordinates": [-114.13992156438134, 51.02930819585162]
		},
		"permitnum": "DP2022-00156",
		"address": "2819 36 ST SW",
		"applicant": "TRICOR DESIGN GROUP",
		"category": "Residential - New Single / Semi / Duplex",
		"description": "NEW: CONTEXTUAL SEMI-DETACHED DWELLING, ACCESSORY RESIDENTIAL BUILDING (GARAGE)",
		"proposedusecode": "C1020; C1368",
		"proposedusedescription": "ACCESSORY RESIDENTIAL BUILDING; CONTEXTUAL SEMI-DETACHED DWELLING",
		"permitteddiscretionary": "Permitted",
		"landusedistrict": "R-CG",
		"landusedistrictdescription": "Residential - Grade-Oriented Infill",
		"statuscurrent": "In Circulation",
		"applieddate": "2022-01-10T00:00:00.000",
		"decisiondate": "2022-02-16T00:00:00.000",
		"releasedate": "2022-04-19T00:00:00.000",
		"mustcommencedate": "2024-02-16T00:00:00.000",
		"decision": "Approval",
		"decisionby": "Development Authority",
		"communitycode": "KIL",
		"communityname": "KILLARNEY/GLENGARRY",
		"ward": "8",
		"quadrant": "SW",
		"latitude": "51.029",
		"longitude": "-114.140",
		"locationcount": "2",
		"locationtypes": "Titled Parcel;Building",
		"locationaddresses": "2819 36 ST SW;2819 36 ST SW",
		"locationsgeojson": "{\"type\":\"MultiPoint\",\"coordinates\":[[-114.1399216,51.0293082],[-114.1400278,51.029353]]}",
		"locationswkt": "MULTIPOINT ((-114.13992156438134 51.02930819585162),(-114.14002777754924 51.029352999468266))"
	}
]
`)

	storedPermits, errS := parseDevelopmentPermits(storedDP)
	fetchedPermits, errF := parseDevelopmentPermits(fetchedDP)

	assert.NoError(t, errS)
	assert.NoError(t, errF)

	createdActions := getDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 0, len(createdActions))
}

func TestContains_DevelopmentPermitAction_StatusUpdate(t *testing.T) {
	storedDP := []byte(`
	[
		{
			"point": {
				"type": "Point",
				"coordinates": [-114.13992156438134, 51.02930819585162]
			},
			"permitnum": "DP2022-00156",
			"address": "2819 36 ST SW",
			"applicant": "TRICOR DESIGN GROUP",
			"category": "Residential - New Single / Semi / Duplex",
			"description": "NEW: CONTEXTUAL SEMI-DETACHED DWELLING, ACCESSORY RESIDENTIAL BUILDING (GARAGE)",
			"proposedusecode": "C1020; C1368",
			"proposedusedescription": "ACCESSORY RESIDENTIAL BUILDING; CONTEXTUAL SEMI-DETACHED DWELLING",
			"permitteddiscretionary": "Permitted",
			"landusedistrict": "R-CG",
			"landusedistrictdescription": "Residential - Grade-Oriented Infill",
			"statuscurrent": "In Circulation",
			"applieddate": "2022-01-10T00:00:00.000",
			"decisiondate": "2022-02-16T00:00:00.000",
			"releasedate": "2022-04-19T00:00:00.000",
			"mustcommencedate": "2024-02-16T00:00:00.000",
			"decision": "Approval",
			"decisionby": "Development Authority",
			"communitycode": "KIL",
			"communityname": "KILLARNEY/GLENGARRY",
			"ward": "8",
			"quadrant": "SW",
			"latitude": "51.029",
			"longitude": "-114.140",
			"locationcount": "2",
			"locationtypes": "Titled Parcel;Building",
			"locationaddresses": "2819 36 ST SW;2819 36 ST SW",
			"locationsgeojson": "{\"type\":\"MultiPoint\",\"coordinates\":[[-114.1399216,51.0293082],[-114.1400278,51.029353]]}",
			"locationswkt": "MULTIPOINT ((-114.13992156438134 51.02930819585162),(-114.14002777754924 51.029352999468266))",
			"github_discussion_id": "D_kwDOJSUT984ATT9k",
			"github_discussion_closed": false
		}
	]
`)
	fetchedDP := []byte(`
[
	{
		"point": {
			"type": "Point",
			"coordinates": [-114.13992156438134, 51.02930819585162]
		},
		"permitnum": "DP2022-00156",
		"address": "2819 36 ST SW",
		"applicant": "TRICOR DESIGN GROUP",
		"category": "Residential - New Single / Semi / Duplex",
		"description": "NEW: CONTEXTUAL SEMI-DETACHED DWELLING, ACCESSORY RESIDENTIAL BUILDING (GARAGE)",
		"proposedusecode": "C1020; C1368",
		"proposedusedescription": "ACCESSORY RESIDENTIAL BUILDING; CONTEXTUAL SEMI-DETACHED DWELLING",
		"permitteddiscretionary": "Permitted",
		"landusedistrict": "R-CG",
		"landusedistrictdescription": "Residential - Grade-Oriented Infill",
		"statuscurrent": "In Progress",
		"applieddate": "2022-01-10T00:00:00.000",
		"decisiondate": "2022-02-16T00:00:00.000",
		"releasedate": "2022-04-19T00:00:00.000",
		"mustcommencedate": "2024-02-16T00:00:00.000",
		"decision": "Approval",
		"decisionby": "Development Authority",
		"communitycode": "KIL",
		"communityname": "KILLARNEY/GLENGARRY",
		"ward": "8",
		"quadrant": "SW",
		"latitude": "51.029",
		"longitude": "-114.140",
		"locationcount": "2",
		"locationtypes": "Titled Parcel;Building",
		"locationaddresses": "2819 36 ST SW;2819 36 ST SW",
		"locationsgeojson": "{\"type\":\"MultiPoint\",\"coordinates\":[[-114.1399216,51.0293082],[-114.1400278,51.029353]]}",
		"locationswkt": "MULTIPOINT ((-114.13992156438134 51.02930819585162),(-114.14002777754924 51.029352999468266))"
	}
]
`)
	expectedActions := []fileaction.FileAction{
		{
			PermitNum: "DP2022-00156",
			Action:    "UPDATE",
		},
	}

	storedPermits, errS := parseDevelopmentPermits(storedDP)
	fetchedPermits, errF := parseDevelopmentPermits(fetchedDP)

	assert.NoError(t, errS)
	assert.NoError(t, errF)

	createdActions := getDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 1, len(createdActions))
	assert.Equal(t, expectedActions[0].PermitNum, createdActions[0].PermitNum)
	assert.Equal(t, expectedActions[0].Action, createdActions[0].Action)
	assert.Contains(t, createdActions[0].Message, "Status")
}

func TestContains_DevelopmentPermitAction_DecisionUpdate(t *testing.T) {
	storedDP := []byte(`
	[
		{
			"point": {
				"type": "Point",
				"coordinates": [-114.13992156438134, 51.02930819585162]
			},
			"permitnum": "DP2022-00156",
			"address": "2819 36 ST SW",
			"applicant": "TRICOR DESIGN GROUP",
			"category": "Residential - New Single / Semi / Duplex",
			"description": "NEW: CONTEXTUAL SEMI-DETACHED DWELLING, ACCESSORY RESIDENTIAL BUILDING (GARAGE)",
			"proposedusecode": "C1020; C1368",
			"proposedusedescription": "ACCESSORY RESIDENTIAL BUILDING; CONTEXTUAL SEMI-DETACHED DWELLING",
			"permitteddiscretionary": "Permitted",
			"landusedistrict": "R-CG",
			"landusedistrictdescription": "Residential - Grade-Oriented Infill",
			"statuscurrent": "In Circulation",
			"applieddate": "2022-01-10T00:00:00.000",
			"decisiondate": "2022-02-16T00:00:00.000",
			"releasedate": "2022-04-19T00:00:00.000",
			"mustcommencedate": "2024-02-16T00:00:00.000",
			"decision": null,
			"decisionby": "Development Authority",
			"communitycode": "KIL",
			"communityname": "KILLARNEY/GLENGARRY",
			"ward": "8",
			"quadrant": "SW",
			"latitude": "51.029",
			"longitude": "-114.140",
			"locationcount": "2",
			"locationtypes": "Titled Parcel;Building",
			"locationaddresses": "2819 36 ST SW;2819 36 ST SW",
			"locationsgeojson": "{\"type\":\"MultiPoint\",\"coordinates\":[[-114.1399216,51.0293082],[-114.1400278,51.029353]]}",
			"locationswkt": "MULTIPOINT ((-114.13992156438134 51.02930819585162),(-114.14002777754924 51.029352999468266))",
			"github_discussion_id": "D_kwDOJSUT984ATT9k",
			"github_discussion_closed": false
		}
	]
`)
	fetchedDP := []byte(`
[
	{
		"point": {
			"type": "Point",
			"coordinates": [-114.13992156438134, 51.02930819585162]
		},
		"permitnum": "DP2022-00156",
		"address": "2819 36 ST SW",
		"applicant": "TRICOR DESIGN GROUP",
		"category": "Residential - New Single / Semi / Duplex",
		"description": "NEW: CONTEXTUAL SEMI-DETACHED DWELLING, ACCESSORY RESIDENTIAL BUILDING (GARAGE)",
		"proposedusecode": "C1020; C1368",
		"proposedusedescription": "ACCESSORY RESIDENTIAL BUILDING; CONTEXTUAL SEMI-DETACHED DWELLING",
		"permitteddiscretionary": "Permitted",
		"landusedistrict": "R-CG",
		"landusedistrictdescription": "Residential - Grade-Oriented Infill",
		"statuscurrent": "In Progress",
		"applieddate": "2022-01-10T00:00:00.000",
		"decisiondate": "2022-02-16T00:00:00.000",
		"releasedate": "2022-04-19T00:00:00.000",
		"mustcommencedate": "2024-02-16T00:00:00.000",
		"decision": "Approval",
		"decisionby": "Development Authority",
		"communitycode": "KIL",
		"communityname": "KILLARNEY/GLENGARRY",
		"ward": "8",
		"quadrant": "SW",
		"latitude": "51.029",
		"longitude": "-114.140",
		"locationcount": "2",
		"locationtypes": "Titled Parcel;Building",
		"locationaddresses": "2819 36 ST SW;2819 36 ST SW",
		"locationsgeojson": "{\"type\":\"MultiPoint\",\"coordinates\":[[-114.1399216,51.0293082],[-114.1400278,51.029353]]}",
		"locationswkt": "MULTIPOINT ((-114.13992156438134 51.02930819585162),(-114.14002777754924 51.029352999468266))"
	}
]
`)
	expectedActions := []fileaction.FileAction{
		{
			PermitNum: "DP2022-00156",
			Action:    "UPDATE",
		},
	}

	storedPermits, errS := parseDevelopmentPermits(storedDP)
	fetchedPermits, errF := parseDevelopmentPermits(fetchedDP)

	assert.NoError(t, errS)
	assert.NoError(t, errF)

	createdActions := getDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 1, len(createdActions))
	assert.Equal(t, expectedActions[0].PermitNum, createdActions[0].PermitNum)
	assert.Equal(t, expectedActions[0].Action, createdActions[0].Action)
	assert.Contains(t, createdActions[0].Message, "Decision")
}

func TestContains_DevelopmentPermitAction_MultiUpdate(t *testing.T) {
	storedDP := []byte(`
	[
		{
			"point": {
				"type": "Point",
				"coordinates": [-114.13992156438134, 51.02930819585162]
			},
			"permitnum": "DP2022-00156",
			"address": "2819 36 ST SW",
			"applicant": "TRICOR DESIGN GROUP",
			"category": "Residential - New Single / Semi / Duplex",
			"description": "NEW: CONTEXTUAL SEMI-DETACHED DWELLING, ACCESSORY RESIDENTIAL BUILDING (GARAGE)",
			"proposedusecode": "C1020; C1368",
			"proposedusedescription": "ACCESSORY RESIDENTIAL BUILDING; CONTEXTUAL SEMI-DETACHED DWELLING",
			"permitteddiscretionary": "Permitted",
			"landusedistrict": "R-CG",
			"landusedistrictdescription": "Residential - Grade-Oriented Infill",
			"statuscurrent": "New",
			"applieddate": "2022-01-10T00:00:00.000",
			"decisiondate": "2022-02-16T00:00:00.000",
			"releasedate": "2022-04-19T00:00:00.000",
			"mustcommencedate": "2024-02-16T00:00:00.000",
			"decision": null,
			"decisionby": "Development Authority",
			"communitycode": "KIL",
			"communityname": "KILLARNEY/GLENGARRY",
			"ward": "8",
			"quadrant": "SW",
			"latitude": "51.029",
			"longitude": "-114.140",
			"locationcount": "2",
			"locationtypes": "Titled Parcel;Building",
			"locationaddresses": "2819 36 ST SW;2819 36 ST SW",
			"locationsgeojson": "{\"type\":\"MultiPoint\",\"coordinates\":[[-114.1399216,51.0293082],[-114.1400278,51.029353]]}",
			"locationswkt": "MULTIPOINT ((-114.13992156438134 51.02930819585162),(-114.14002777754924 51.029352999468266))",
			"github_discussion_id": "D_kwDOJSUT984ATT9k",
			"github_discussion_closed": false
		}
	]
`)
	fetchedDP := []byte(`
[
	{
		"point": {
			"type": "Point",
			"coordinates": [-114.13992156438134, 51.02930819585162]
		},
		"permitnum": "DP2022-00156",
		"address": "2819 36 ST SW",
		"applicant": "TRICOR DESIGN GROUP",
		"category": "Residential - New Single / Semi / Duplex",
		"description": "NEW: CONTEXTUAL SEMI-DETACHED DWELLING, ACCESSORY RESIDENTIAL BUILDING (GARAGE)",
		"proposedusecode": "C1020; C1368",
		"proposedusedescription": "ACCESSORY RESIDENTIAL BUILDING; CONTEXTUAL SEMI-DETACHED DWELLING",
		"permitteddiscretionary": "Permitted",
		"landusedistrict": "R-CG",
		"landusedistrictdescription": "Residential - Grade-Oriented Infill",
		"statuscurrent": "In Progress",
		"applieddate": "2022-01-10T00:00:00.000",
		"decisiondate": "2022-02-16T00:00:00.000",
		"releasedate": "2022-04-19T00:00:00.000",
		"mustcommencedate": "2024-02-16T00:00:00.000",
		"decision": "Approval",
		"decisionby": "Development Authority",
		"communitycode": "KIL",
		"communityname": "KILLARNEY/GLENGARRY",
		"ward": "8",
		"quadrant": "SW",
		"latitude": "51.029",
		"longitude": "-114.140",
		"locationcount": "2",
		"locationtypes": "Titled Parcel;Building",
		"locationaddresses": "2819 36 ST SW;2819 36 ST SW",
		"locationsgeojson": "{\"type\":\"MultiPoint\",\"coordinates\":[[-114.1399216,51.0293082],[-114.1400278,51.029353]]}",
		"locationswkt": "MULTIPOINT ((-114.13992156438134 51.02930819585162),(-114.14002777754924 51.029352999468266))"
	}
]
`)
	expectedActions := []fileaction.FileAction{
		{
			PermitNum: "DP2022-00156",
			Action:    "UPDATE",
		},
	}

	storedPermits, errS := parseDevelopmentPermits(storedDP)
	fetchedPermits, errF := parseDevelopmentPermits(fetchedDP)

	assert.NoError(t, errS)
	assert.NoError(t, errF)

	createdActions := getDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 1, len(createdActions))
	assert.Equal(t, expectedActions[0].PermitNum, createdActions[0].PermitNum)
	assert.Equal(t, expectedActions[0].Action, createdActions[0].Action)
	assert.Contains(t, createdActions[0].Message, "Decision")
	assert.Contains(t, createdActions[0].Message, "Status")
}

func TestContains_DevelopmentPermitAction_StatusUpdateReleased(t *testing.T) {
	storedDP := []byte(`
	[
		{
			"point": {
				"type": "Point",
				"coordinates": [-114.13992156438134, 51.02930819585162]
			},
			"permitnum": "DP2022-00156",
			"address": "2819 36 ST SW",
			"applicant": "TRICOR DESIGN GROUP",
			"category": "Residential - New Single / Semi / Duplex",
			"description": "NEW: CONTEXTUAL SEMI-DETACHED DWELLING, ACCESSORY RESIDENTIAL BUILDING (GARAGE)",
			"proposedusecode": "C1020; C1368",
			"proposedusedescription": "ACCESSORY RESIDENTIAL BUILDING; CONTEXTUAL SEMI-DETACHED DWELLING",
			"permitteddiscretionary": "Permitted",
			"landusedistrict": "R-CG",
			"landusedistrictdescription": "Residential - Grade-Oriented Infill",
			"statuscurrent": "In Circulation",
			"applieddate": "2022-01-10T00:00:00.000",
			"decisiondate": "2022-02-16T00:00:00.000",
			"releasedate": "2022-04-19T00:00:00.000",
			"mustcommencedate": "2024-02-16T00:00:00.000",
			"decision": null,
			"decisionby": "Development Authority",
			"communitycode": "KIL",
			"communityname": "KILLARNEY/GLENGARRY",
			"ward": "8",
			"quadrant": "SW",
			"latitude": "51.029",
			"longitude": "-114.140",
			"locationcount": "2",
			"locationtypes": "Titled Parcel;Building",
			"locationaddresses": "2819 36 ST SW;2819 36 ST SW",
			"locationsgeojson": "{\"type\":\"MultiPoint\",\"coordinates\":[[-114.1399216,51.0293082],[-114.1400278,51.029353]]}",
			"locationswkt": "MULTIPOINT ((-114.13992156438134 51.02930819585162),(-114.14002777754924 51.029352999468266))",
			"github_discussion_id": "D_kwDOJSUT984ATT9k",
			"github_discussion_closed": false
		}
	]
`)
	fetchedDP := []byte(`
[
	{
		"point": {
			"type": "Point",
			"coordinates": [-114.13992156438134, 51.02930819585162]
		},
		"permitnum": "DP2022-00156",
		"address": "2819 36 ST SW",
		"applicant": "TRICOR DESIGN GROUP",
		"category": "Residential - New Single / Semi / Duplex",
		"description": "NEW: CONTEXTUAL SEMI-DETACHED DWELLING, ACCESSORY RESIDENTIAL BUILDING (GARAGE)",
		"proposedusecode": "C1020; C1368",
		"proposedusedescription": "ACCESSORY RESIDENTIAL BUILDING; CONTEXTUAL SEMI-DETACHED DWELLING",
		"permitteddiscretionary": "Permitted",
		"landusedistrict": "R-CG",
		"landusedistrictdescription": "Residential - Grade-Oriented Infill",
		"statuscurrent": "Released",
		"applieddate": "2022-01-10T00:00:00.000",
		"decisiondate": "2022-02-16T00:00:00.000",
		"releasedate": "2022-04-19T00:00:00.000",
		"mustcommencedate": "2024-02-16T00:00:00.000",
		"decision": "Approval",
		"decisionby": "Development Authority",
		"communitycode": "KIL",
		"communityname": "KILLARNEY/GLENGARRY",
		"ward": "8",
		"quadrant": "SW",
		"latitude": "51.029",
		"longitude": "-114.140",
		"locationcount": "2",
		"locationtypes": "Titled Parcel;Building",
		"locationaddresses": "2819 36 ST SW;2819 36 ST SW",
		"locationsgeojson": "{\"type\":\"MultiPoint\",\"coordinates\":[[-114.1399216,51.0293082],[-114.1400278,51.029353]]}",
		"locationswkt": "MULTIPOINT ((-114.13992156438134 51.02930819585162),(-114.14002777754924 51.029352999468266))"
	}
]
`)
	expectedActions := []fileaction.FileAction{
		{
			PermitNum: "DP2022-00156",
			Action:    "CLOSE",
		},
	}

	storedPermits, errS := parseDevelopmentPermits(storedDP)
	fetchedPermits, errF := parseDevelopmentPermits(fetchedDP)

	assert.NoError(t, errS)
	assert.NoError(t, errF)

	createdActions := getDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 1, len(createdActions))
	assert.Equal(t, expectedActions[0].PermitNum, createdActions[0].PermitNum)
	assert.Equal(t, expectedActions[0].Action, createdActions[0].Action)
	assert.Contains(t, createdActions[0].Message, "Decision")
	assert.Contains(t, createdActions[0].Message, "Status")
	assert.Contains(t, createdActions[0].Message, "Closing")
}

func TestContains_DevelopmentPermitAction_Create(t *testing.T) {
	storedDP := []byte(`[]`)
	fetchedDP := []byte(`
[
	{
		"point": {
			"type": "Point",
			"coordinates": [-114.13992156438134, 51.02930819585162]
		},
		"permitnum": "DP2022-00156",
		"address": "2819 36 ST SW",
		"applicant": "TRICOR DESIGN GROUP",
		"category": "Residential - New Single / Semi / Duplex",
		"description": "NEW: CONTEXTUAL SEMI-DETACHED DWELLING, ACCESSORY RESIDENTIAL BUILDING (GARAGE)",
		"proposedusecode": "C1020; C1368",
		"proposedusedescription": "ACCESSORY RESIDENTIAL BUILDING; CONTEXTUAL SEMI-DETACHED DWELLING",
		"permitteddiscretionary": "Permitted",
		"landusedistrict": "R-CG",
		"landusedistrictdescription": "Residential - Grade-Oriented Infill",
		"statuscurrent": "In Progress",
		"applieddate": "2022-01-10T00:00:00.000",
		"decisiondate": "2022-02-16T00:00:00.000",
		"releasedate": "2022-04-19T00:00:00.000",
		"mustcommencedate": "2024-02-16T00:00:00.000",
		"decision": "Approval",
		"decisionby": "Development Authority",
		"communitycode": "KIL",
		"communityname": "KILLARNEY/GLENGARRY",
		"ward": "8",
		"quadrant": "SW",
		"latitude": "51.029",
		"longitude": "-114.140",
		"locationcount": "2",
		"locationtypes": "Titled Parcel;Building",
		"locationaddresses": "2819 36 ST SW;2819 36 ST SW",
		"locationsgeojson": "{\"type\":\"MultiPoint\",\"coordinates\":[[-114.1399216,51.0293082],[-114.1400278,51.029353]]}",
		"locationswkt": "MULTIPOINT ((-114.13992156438134 51.02930819585162),(-114.14002777754924 51.029352999468266))"
	}
]
`)
	expectedActions := []fileaction.FileAction{
		{
			PermitNum: "DP2022-00156",
			Action:    "CREATE",
		},
	}

	storedPermits, errS := parseDevelopmentPermits(storedDP)
	fetchedPermits, errF := parseDevelopmentPermits(fetchedDP)

	assert.NoError(t, errS)
	assert.NoError(t, errF)

	createdActions := getDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 1, len(createdActions))
	assert.Equal(t, expectedActions[0].PermitNum, createdActions[0].PermitNum)
	assert.Equal(t, expectedActions[0].Action, createdActions[0].Action)
}

func TestDevelopmentPermitActions_CreateVsUpdate(t *testing.T) {
	// Test CREATE action - new permit not in stored data
	storedDP := []byte(`[]`)
	fetchedDP := []byte(`
[
	{
		"point": {
			"type": "Point",
			"coordinates": [-114.13992156438134, 51.02930819585162]
		},
		"permitnum": "DP2025-12345",
		"address": "123 Test ST SW",
		"statuscurrent": "In Progress",
		"applieddate": "2025-01-01T00:00:00.000"
	}
]
`)

	storedPermits, errS := parseDevelopmentPermits(storedDP)
	fetchedPermits, errF := parseDevelopmentPermits(fetchedDP)

	assert.NoError(t, errS)
	assert.NoError(t, errF)

	actions := getDevelopmentPermitActions(fetchedPermits, storedPermits)

	// Should generate CREATE action for new permit
	assert.Equal(t, 1, len(actions))
	assert.Equal(t, "DP2025-12345", actions[0].PermitNum)
	assert.Equal(t, "CREATE", actions[0].Action)
}

func TestDevelopmentPermitActions_UpdateAction(t *testing.T) {
	// Test UPDATE action - existing permit with status change
	storedDP := []byte(`
[
	{
		"point": {
			"type": "Point",
			"coordinates": [-114.13992156438134, 51.02930819585162]
		},
		"permitnum": "DP2025-12345",
		"address": "123 Test ST SW",
		"statuscurrent": "In Progress",
		"applieddate": "2025-01-01T00:00:00.000"
	}
]
`)
	fetchedDP := []byte(`
[
	{
		"point": {
			"type": "Point",
			"coordinates": [-114.13992156438134, 51.02930819585162]
		},
		"permitnum": "DP2025-12345",
		"address": "123 Test ST SW",
		"statuscurrent": "Released",
		"applieddate": "2025-01-01T00:00:00.000",
		"decision": "Approval",
		"decisiondate": "2025-01-15T00:00:00.000"
	}
]
`)

	storedPermits, errS := parseDevelopmentPermits(storedDP)
	fetchedPermits, errF := parseDevelopmentPermits(fetchedDP)

	assert.NoError(t, errS)
	assert.NoError(t, errF)

	actions := getDevelopmentPermitActions(fetchedPermits, storedPermits)

	// Should generate CLOSE action for status change to Released
	assert.Equal(t, 1, len(actions))
	assert.Equal(t, "DP2025-12345", actions[0].PermitNum)
	assert.Equal(t, "CLOSE", actions[0].Action)
	assert.Contains(t, actions[0].Message, "Status updated")
	assert.Contains(t, actions[0].Message, "Closing file")
}

func TestDevelopmentPermitActions_NoActionWhenAlreadyClosed(t *testing.T) {
	// Test no action when permit is already in closed status
	storedDP := []byte(`
[
	{
		"point": {
			"type": "Point",
			"coordinates": [-114.13992156438134, 51.02930819585162]
		},
		"permitnum": "DP2025-12345",
		"address": "123 Test ST SW",
		"statuscurrent": "Released",
		"applieddate": "2025-01-01T00:00:00.000",
		"decision": "Approval",
		"decisiondate": "2025-01-15T00:00:00.000"
	}
]
`)
	fetchedDP := []byte(`
[
	{
		"point": {
			"type": "Point",
			"coordinates": [-114.13992156438134, 51.02930819585162]
		},
		"permitnum": "DP2025-12345",
		"address": "123 Test ST SW",
		"statuscurrent": "Released",
		"applieddate": "2025-01-01T00:00:00.000",
		"decision": "Approval",
		"decisiondate": "2025-01-15T00:00:00.000"
	}
]
`)

	storedPermits, errS := parseDevelopmentPermits(storedDP)
	fetchedPermits, errF := parseDevelopmentPermits(fetchedDP)

	assert.NoError(t, errS)
	assert.NoError(t, errF)

	actions := getDevelopmentPermitActions(fetchedPermits, storedPermits)

	// Should generate no actions since permit is already closed and unchanged
	assert.Equal(t, 0, len(actions))
}

func TestGetMostRecentTimestamp_WithStateHistory(t *testing.T) {
	dp := DevelopmentPermit{
		PermitNum:    "DP2025-12345",
		AppliedDate:  strPtr("2025-01-01T10:00:00.000"),
		DecisionDate: strPtr("2025-01-15T14:30:00.000"),
		StateHistory: []StateChange{
			{
				Status:    "in progress",
				Timestamp: "2025-01-01T10:00:00-07:00",
			},
			{
				Status:    "released",
				Timestamp: "2025-01-20T16:45:00-07:00", // This should be the most recent
			},
		},
	}

	mostRecent := dp.getMostRecentTimestamp()

	// Should use the most recent state history timestamp
	expected, _ := time.Parse(time.RFC3339, "2025-01-20T16:45:00-07:00")
	assert.Equal(t, expected.Unix(), mostRecent.Unix())
}

func TestGetMostRecentTimestamp_WithoutStateHistory(t *testing.T) {
	dp := DevelopmentPermit{
		PermitNum:    "DP2025-12345",
		AppliedDate:  strPtr("2025-01-01T10:00:00.000"),
		DecisionDate: strPtr("2025-01-15T14:30:00.000"),
		ReleaseDate:  strPtr("2025-01-20T16:45:00.000"), // This should be the most recent
		StateHistory: []StateChange{},                   // Empty state history
	}

	mostRecent := dp.getMostRecentTimestamp()

	// Should use the release date as the most recent
	expected, _ := time.Parse("2006-01-02T15:04:05.000", "2025-01-20T16:45:00.000")
	assert.Equal(t, expected.Unix(), mostRecent.Unix())
}

func TestGetMostRecentTimestamp_NoTimestamps(t *testing.T) {
	dp := DevelopmentPermit{
		PermitNum:    "DP2025-12345",
		StateHistory: []StateChange{},
	}

	mostRecent := dp.getMostRecentTimestamp()

	// Should use current time when no timestamps are available
	now := time.Now()
	// Allow for a small time difference since the function calls time.Now()
	assert.WithinDuration(t, now, mostRecent, time.Second)
}

// Helper function for string pointers
func strPtr(s string) *string {
	return &s
}
