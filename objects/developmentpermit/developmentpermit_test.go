package developmentpermit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FindDevelopmentPermit_ItemFound(t *testing.T) {
	dp1 := DevelopmentPermit{PermitNum: "123"}
	dp2 := DevelopmentPermit{PermitNum: "456"}
	dp3 := DevelopmentPermit{PermitNum: "789"}
	searchSlice := []DevelopmentPermit{dp1, dp2, dp3}

	result := FindDevelopmentPermit(searchSlice, "456")

	assert.NotNil(t, result)
	assert.Equal(t, &dp2, result)
}

func Test_FindDevelopmentPermit_ItemNotFound(t *testing.T) {
	dp1 := DevelopmentPermit{PermitNum: "123"}
	//	dp2 := DevelopmentPermit{PermitNum: "456"}
	dp3 := DevelopmentPermit{PermitNum: "789"}
	searchSlice := []DevelopmentPermit{dp1, dp3}

	result := FindDevelopmentPermit(searchSlice, "456")

	assert.Nil(t, result)
}

func Test_FindDevelopmentPermit_EmptySlice(t *testing.T) {
	searchSlice := []DevelopmentPermit{}

	result := FindDevelopmentPermit(searchSlice, "456")

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
	storedPermit := DevelopmentPermit{
		PermitNum:              "DP2023-TEST",
		Address:                "123 Testing Street",
		Applicant:              "Unit Test Inc",
		Category:               "Residential - Single Family",
		Description:            "NEW:RESIDENTIAL DEVELOPMENT (1 BUILDING), ACCESSORY RESIDENTIAL BUILDING (GARAGE) ",
		ProposedUseCode:        "T0001; T0002",
		PermittedDiscretion:    "Permitted",
		LandUseDistrict:        "R-C2",
		LandUseDistrictDesc:    "Residential - Contextual One/Two Dwelling",
		StatusCurrent:          "Under Review",
		AppliedDate:            "2023-04-08T00:00:00.000",
		CommunityCode:          "KIL",
		CommunityName:          "KILLARNEY/GLENGARRY",
		Ward:                   "8",
		Quadrant:               "SW",
		Latitude:               "51.034",
		Longitude:              "-114.122",
		LocationCount:          "1",
		LocationTypes:          "Titled Parcel",
		LocationAddresses:      "123 Testing Street",
		LocationsGeoJSON:       "{\"type\":\"Point\",\"coordinates\":[-114.1224605,51.0338002]}",
		LocationsWKT:           "POINT (-114.12246048560655 51.03380020179407)",
		DecisionDate:           "",
		MustCommenceDate:       "",
		Decision:               "",
		DecisionBy:             "",
		ReleaseDate:            "",
		GithubDiscussionClosed: false,
	}
	fetchedPermit := DevelopmentPermit{
		PermitNum:              "DP2023-TEST",
		Address:                "123 Testing Street",
		Applicant:              "Unit Test Inc",
		Category:               "Residential - Single Family",
		Description:            "NEW:RESIDENTIAL DEVELOPMENT (1 BUILDING), ACCESSORY RESIDENTIAL BUILDING (GARAGE) ",
		ProposedUseCode:        "T0001; T0002",
		PermittedDiscretion:    "Permitted",
		LandUseDistrict:        "R-C2",
		LandUseDistrictDesc:    "Residential - Contextual One/Two Dwelling",
		StatusCurrent:          "Under Review",
		AppliedDate:            "2023-04-08T00:00:00.000",
		CommunityCode:          "KIL",
		CommunityName:          "KILLARNEY/GLENGARRY",
		Ward:                   "8",
		Quadrant:               "SW",
		Latitude:               "51.034",
		Longitude:              "-114.122",
		LocationCount:          "1",
		LocationTypes:          "Titled Parcel",
		LocationAddresses:      "123 Testing Street",
		LocationsGeoJSON:       "{\"type\":\"Point\",\"coordinates\":[-114.1224605,51.0338002]}",
		LocationsWKT:           "POINT (-114.12246048560655 51.03380020179407)",
		DecisionDate:           "",
		MustCommenceDate:       "",
		Decision:               "",
		DecisionBy:             "",
		ReleaseDate:            "",
		GithubDiscussionClosed: false,
	}

	fetchedPermits := []DevelopmentPermit{fetchedPermit}
	storedPermits := []DevelopmentPermit{storedPermit}

	actions := GetDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 0, len(actions))
}

func TestContains_DevelopmentPermitAction_StatusUpdate(t *testing.T) {
	storedPermit := DevelopmentPermit{
		PermitNum:              "DP2023-TEST",
		Address:                "123 Testing Street",
		Applicant:              "Unit Test Inc",
		Category:               "Residential - Single Family",
		Description:            "NEW:RESIDENTIAL DEVELOPMENT (1 BUILDING), ACCESSORY RESIDENTIAL BUILDING (GARAGE) ",
		ProposedUseCode:        "T0001; T0002",
		PermittedDiscretion:    "Permitted",
		LandUseDistrict:        "R-C2",
		LandUseDistrictDesc:    "Residential - Contextual One/Two Dwelling",
		StatusCurrent:          "Under Review",
		AppliedDate:            "2023-04-08T00:00:00.000",
		CommunityCode:          "KIL",
		CommunityName:          "KILLARNEY/GLENGARRY",
		Ward:                   "8",
		Quadrant:               "SW",
		Latitude:               "51.034",
		Longitude:              "-114.122",
		LocationCount:          "1",
		LocationTypes:          "Titled Parcel",
		LocationAddresses:      "123 Testing Street",
		LocationsGeoJSON:       "{\"type\":\"Point\",\"coordinates\":[-114.1224605,51.0338002]}",
		LocationsWKT:           "POINT (-114.12246048560655 51.03380020179407)",
		DecisionDate:           "",
		MustCommenceDate:       "",
		Decision:               "",
		DecisionBy:             "",
		ReleaseDate:            "",
		GithubDiscussionClosed: false,
	}
	fetchedPermit := DevelopmentPermit{
		PermitNum:              "DP2023-TEST",
		Address:                "123 Testing Street",
		Applicant:              "Unit Test Inc",
		Category:               "Residential - Single Family",
		Description:            "NEW:RESIDENTIAL DEVELOPMENT (1 BUILDING), ACCESSORY RESIDENTIAL BUILDING (GARAGE) ",
		ProposedUseCode:        "T0001; T0002",
		PermittedDiscretion:    "Permitted",
		LandUseDistrict:        "R-C2",
		LandUseDistrictDesc:    "Residential - Contextual One/Two Dwelling",
		StatusCurrent:          "In circulation",
		AppliedDate:            "2023-04-08T00:00:00.000",
		CommunityCode:          "KIL",
		CommunityName:          "KILLARNEY/GLENGARRY",
		Ward:                   "8",
		Quadrant:               "SW",
		Latitude:               "51.034",
		Longitude:              "-114.122",
		LocationCount:          "1",
		LocationTypes:          "Titled Parcel",
		LocationAddresses:      "123 Testing Street",
		LocationsGeoJSON:       "{\"type\":\"Point\",\"coordinates\":[-114.1224605,51.0338002]}",
		LocationsWKT:           "POINT (-114.12246048560655 51.03380020179407)",
		DecisionDate:           "",
		MustCommenceDate:       "",
		Decision:               "",
		DecisionBy:             "",
		ReleaseDate:            "",
		GithubDiscussionClosed: false,
	}

	fetchedPermits := []DevelopmentPermit{fetchedPermit}
	storedPermits := []DevelopmentPermit{storedPermit}

	actions := GetDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 1, len(actions))
	assert.Equal(t, "DP2023-TEST", actions[0].PermitNum)
	assert.Contains(t, actions[0].Action, "UPDATE")
	assert.Contains(t, actions[0].Message, "Status")
}

func TestContains_DevelopmentPermitAction_DecisionUpdate(t *testing.T) {
	storedPermit := DevelopmentPermit{
		PermitNum:              "DP2023-TEST",
		Address:                "123 Testing Street",
		Applicant:              "Unit Test Inc",
		Category:               "Residential - Single Family",
		Description:            "NEW:RESIDENTIAL DEVELOPMENT (1 BUILDING), ACCESSORY RESIDENTIAL BUILDING (GARAGE) ",
		ProposedUseCode:        "T0001; T0002",
		PermittedDiscretion:    "Permitted",
		LandUseDistrict:        "R-C2",
		LandUseDistrictDesc:    "Residential - Contextual One/Two Dwelling",
		StatusCurrent:          "Under Review",
		AppliedDate:            "2023-04-08T00:00:00.000",
		CommunityCode:          "KIL",
		CommunityName:          "KILLARNEY/GLENGARRY",
		Ward:                   "8",
		Quadrant:               "SW",
		Latitude:               "51.034",
		Longitude:              "-114.122",
		LocationCount:          "1",
		LocationTypes:          "Titled Parcel",
		LocationAddresses:      "123 Testing Street",
		LocationsGeoJSON:       "{\"type\":\"Point\",\"coordinates\":[-114.1224605,51.0338002]}",
		LocationsWKT:           "POINT (-114.12246048560655 51.03380020179407)",
		DecisionDate:           "",
		MustCommenceDate:       "",
		Decision:               "",
		DecisionBy:             "",
		ReleaseDate:            "",
		GithubDiscussionClosed: false,
	}
	fetchedPermit := DevelopmentPermit{
		PermitNum:              "DP2023-TEST",
		Address:                "123 Testing Street",
		Applicant:              "Unit Test Inc",
		Category:               "Residential - Single Family",
		Description:            "NEW:RESIDENTIAL DEVELOPMENT (1 BUILDING), ACCESSORY RESIDENTIAL BUILDING (GARAGE) ",
		ProposedUseCode:        "T0001; T0002",
		PermittedDiscretion:    "Permitted",
		LandUseDistrict:        "R-C2",
		LandUseDistrictDesc:    "Residential - Contextual One/Two Dwelling",
		StatusCurrent:          "Under Review",
		AppliedDate:            "2023-04-08T00:00:00.000",
		CommunityCode:          "KIL",
		CommunityName:          "KILLARNEY/GLENGARRY",
		Ward:                   "8",
		Quadrant:               "SW",
		Latitude:               "51.034",
		Longitude:              "-114.122",
		LocationCount:          "1",
		LocationTypes:          "Titled Parcel",
		LocationAddresses:      "123 Testing Street",
		LocationsGeoJSON:       "{\"type\":\"Point\",\"coordinates\":[-114.1224605,51.0338002]}",
		LocationsWKT:           "POINT (-114.12246048560655 51.03380020179407)",
		DecisionDate:           "",
		MustCommenceDate:       "",
		Decision:               "Approved",
		DecisionBy:             "",
		ReleaseDate:            "",
		GithubDiscussionClosed: false,
	}

	fetchedPermits := []DevelopmentPermit{fetchedPermit}
	storedPermits := []DevelopmentPermit{storedPermit}

	actions := GetDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 1, len(actions))
	assert.Equal(t, "DP2023-TEST", actions[0].PermitNum)
	assert.Contains(t, actions[0].Action, "UPDATE")
	assert.Contains(t, actions[0].Message, "Decision")
}

func TestContains_DevelopmentPermitAction_MustCommenceByUpdate(t *testing.T) {
	storedPermit := DevelopmentPermit{
		PermitNum:              "DP2023-TEST",
		Address:                "123 Testing Street",
		Applicant:              "Unit Test Inc",
		Category:               "Residential - Single Family",
		Description:            "NEW:RESIDENTIAL DEVELOPMENT (1 BUILDING), ACCESSORY RESIDENTIAL BUILDING (GARAGE) ",
		ProposedUseCode:        "T0001; T0002",
		PermittedDiscretion:    "Permitted",
		LandUseDistrict:        "R-C2",
		LandUseDistrictDesc:    "Residential - Contextual One/Two Dwelling",
		StatusCurrent:          "Under Review",
		AppliedDate:            "2023-04-08T00:00:00.000",
		CommunityCode:          "KIL",
		CommunityName:          "KILLARNEY/GLENGARRY",
		Ward:                   "8",
		Quadrant:               "SW",
		Latitude:               "51.034",
		Longitude:              "-114.122",
		LocationCount:          "1",
		LocationTypes:          "Titled Parcel",
		LocationAddresses:      "123 Testing Street",
		LocationsGeoJSON:       "{\"type\":\"Point\",\"coordinates\":[-114.1224605,51.0338002]}",
		LocationsWKT:           "POINT (-114.12246048560655 51.03380020179407)",
		DecisionDate:           "",
		MustCommenceDate:       "",
		Decision:               "",
		DecisionBy:             "",
		ReleaseDate:            "",
		GithubDiscussionClosed: false,
	}
	fetchedPermit := DevelopmentPermit{
		PermitNum:              "DP2023-TEST",
		Address:                "123 Testing Street",
		Applicant:              "Unit Test Inc",
		Category:               "Residential - Single Family",
		Description:            "NEW:RESIDENTIAL DEVELOPMENT (1 BUILDING), ACCESSORY RESIDENTIAL BUILDING (GARAGE) ",
		ProposedUseCode:        "T0001; T0002",
		PermittedDiscretion:    "Permitted",
		LandUseDistrict:        "R-C2",
		LandUseDistrictDesc:    "Residential - Contextual One/Two Dwelling",
		StatusCurrent:          "Under Review",
		AppliedDate:            "2023-04-08T00:00:00.000",
		CommunityCode:          "KIL",
		CommunityName:          "KILLARNEY/GLENGARRY",
		Ward:                   "8",
		Quadrant:               "SW",
		Latitude:               "51.034",
		Longitude:              "-114.122",
		LocationCount:          "1",
		LocationTypes:          "Titled Parcel",
		LocationAddresses:      "123 Testing Street",
		LocationsGeoJSON:       "{\"type\":\"Point\",\"coordinates\":[-114.1224605,51.0338002]}",
		LocationsWKT:           "POINT (-114.12246048560655 51.03380020179407)",
		DecisionDate:           "",
		MustCommenceDate:       "2023-05-01T00:00:00.000",
		Decision:               "",
		DecisionBy:             "",
		ReleaseDate:            "",
		GithubDiscussionClosed: false,
	}

	fetchedPermits := []DevelopmentPermit{fetchedPermit}
	storedPermits := []DevelopmentPermit{storedPermit}

	actions := GetDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 1, len(actions))
	assert.Equal(t, "DP2023-TEST", actions[0].PermitNum)
	assert.Contains(t, actions[0].Action, "UPDATE")
	assert.Contains(t, actions[0].Message, "Must Commence")
}

func TestContains_DevelopmentPermitAction_MultipleUpdates(t *testing.T) {
	storedPermit := DevelopmentPermit{
		PermitNum:              "DP2023-TEST",
		Address:                "123 Testing Street",
		Applicant:              "Unit Test Inc",
		Category:               "Residential - Single Family",
		Description:            "NEW:RESIDENTIAL DEVELOPMENT (1 BUILDING), ACCESSORY RESIDENTIAL BUILDING (GARAGE) ",
		ProposedUseCode:        "T0001; T0002",
		PermittedDiscretion:    "Permitted",
		LandUseDistrict:        "R-C2",
		LandUseDistrictDesc:    "Residential - Contextual One/Two Dwelling",
		StatusCurrent:          "Under Review",
		AppliedDate:            "2023-04-08T00:00:00.000",
		CommunityCode:          "KIL",
		CommunityName:          "KILLARNEY/GLENGARRY",
		Ward:                   "8",
		Quadrant:               "SW",
		Latitude:               "51.034",
		Longitude:              "-114.122",
		LocationCount:          "1",
		LocationTypes:          "Titled Parcel",
		LocationAddresses:      "123 Testing Street",
		LocationsGeoJSON:       "{\"type\":\"Point\",\"coordinates\":[-114.1224605,51.0338002]}",
		LocationsWKT:           "POINT (-114.12246048560655 51.03380020179407)",
		DecisionDate:           "",
		MustCommenceDate:       "",
		Decision:               "",
		DecisionBy:             "",
		ReleaseDate:            "",
		GithubDiscussionClosed: false,
	}
	fetchedPermit := DevelopmentPermit{
		PermitNum:              "DP2023-TEST",
		Address:                "123 Testing Street",
		Applicant:              "Unit Test Inc",
		Category:               "Residential - Single Family",
		Description:            "NEW:RESIDENTIAL DEVELOPMENT (1 BUILDING), ACCESSORY RESIDENTIAL BUILDING (GARAGE) ",
		ProposedUseCode:        "T0001; T0002",
		PermittedDiscretion:    "Permitted",
		LandUseDistrict:        "R-C2",
		LandUseDistrictDesc:    "Residential - Contextual One/Two Dwelling",
		StatusCurrent:          "Approved",
		AppliedDate:            "2023-04-08T00:00:00.000",
		CommunityCode:          "KIL",
		CommunityName:          "KILLARNEY/GLENGARRY",
		Ward:                   "8",
		Quadrant:               "SW",
		Latitude:               "51.034",
		Longitude:              "-114.122",
		LocationCount:          "1",
		LocationTypes:          "Titled Parcel",
		LocationAddresses:      "123 Testing Street",
		LocationsGeoJSON:       "{\"type\":\"Point\",\"coordinates\":[-114.1224605,51.0338002]}",
		LocationsWKT:           "POINT (-114.12246048560655 51.03380020179407)",
		DecisionDate:           "2023-04-08T00:00:00.000",
		MustCommenceDate:       "2023-05-01T00:00:00.000",
		Decision:               "Approved",
		DecisionBy:             "",
		ReleaseDate:            "",
		GithubDiscussionClosed: false,
	}

	fetchedPermits := []DevelopmentPermit{fetchedPermit}
	storedPermits := []DevelopmentPermit{storedPermit}

	actions := GetDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 1, len(actions))
	assert.Equal(t, "DP2023-TEST", actions[0].PermitNum)
	assert.Contains(t, actions[0].Action, "UPDATE")
	assert.Contains(t, actions[0].Message, "Status")
	assert.Contains(t, actions[0].Message, "Decision")
	assert.Contains(t, actions[0].Message, "Must Commence")
}

func TestContains_DevelopmentPermitAction_StatusUpdateReleased(t *testing.T) {
	storedPermit := DevelopmentPermit{
		PermitNum:              "DP2023-TEST",
		Address:                "123 Testing Street",
		Applicant:              "Unit Test Inc",
		Category:               "Residential - Single Family",
		Description:            "NEW:RESIDENTIAL DEVELOPMENT (1 BUILDING), ACCESSORY RESIDENTIAL BUILDING (GARAGE) ",
		ProposedUseCode:        "T0001; T0002",
		PermittedDiscretion:    "Permitted",
		LandUseDistrict:        "R-C2",
		LandUseDistrictDesc:    "Residential - Contextual One/Two Dwelling",
		StatusCurrent:          "Under Review",
		AppliedDate:            "2023-04-08T00:00:00.000",
		CommunityCode:          "KIL",
		CommunityName:          "KILLARNEY/GLENGARRY",
		Ward:                   "8",
		Quadrant:               "SW",
		Latitude:               "51.034",
		Longitude:              "-114.122",
		LocationCount:          "1",
		LocationTypes:          "Titled Parcel",
		LocationAddresses:      "123 Testing Street",
		LocationsGeoJSON:       "{\"type\":\"Point\",\"coordinates\":[-114.1224605,51.0338002]}",
		LocationsWKT:           "POINT (-114.12246048560655 51.03380020179407)",
		DecisionDate:           "",
		MustCommenceDate:       "",
		Decision:               "",
		DecisionBy:             "",
		ReleaseDate:            "",
		GithubDiscussionClosed: false,
	}
	fetchedPermit := DevelopmentPermit{
		PermitNum:              "DP2023-TEST",
		Address:                "123 Testing Street",
		Applicant:              "Unit Test Inc",
		Category:               "Residential - Single Family",
		Description:            "NEW:RESIDENTIAL DEVELOPMENT (1 BUILDING), ACCESSORY RESIDENTIAL BUILDING (GARAGE) ",
		ProposedUseCode:        "T0001; T0002",
		PermittedDiscretion:    "Permitted",
		LandUseDistrict:        "R-C2",
		LandUseDistrictDesc:    "Residential - Contextual One/Two Dwelling",
		StatusCurrent:          "Released",
		AppliedDate:            "2023-04-08T00:00:00.000",
		CommunityCode:          "KIL",
		CommunityName:          "KILLARNEY/GLENGARRY",
		Ward:                   "8",
		Quadrant:               "SW",
		Latitude:               "51.034",
		Longitude:              "-114.122",
		LocationCount:          "1",
		LocationTypes:          "Titled Parcel",
		LocationAddresses:      "123 Testing Street",
		LocationsGeoJSON:       "{\"type\":\"Point\",\"coordinates\":[-114.1224605,51.0338002]}",
		LocationsWKT:           "POINT (-114.12246048560655 51.03380020179407)",
		DecisionDate:           "",
		MustCommenceDate:       "",
		Decision:               "",
		DecisionBy:             "",
		ReleaseDate:            "",
		GithubDiscussionClosed: false,
	}

	fetchedPermits := []DevelopmentPermit{fetchedPermit}
	storedPermits := []DevelopmentPermit{storedPermit}

	actions := GetDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 1, len(actions))
	assert.Equal(t, "DP2023-TEST", actions[0].PermitNum)
	assert.Contains(t, actions[0].Action, "CLOSE")
	assert.Contains(t, actions[0].Message, "Status")
	assert.Contains(t, actions[0].Message, "Closing")
}

func TestContains_DevelopmentPermitAction_DiscussionClosed(t *testing.T) {
	storedPermit := DevelopmentPermit{
		PermitNum:              "DP2023-TEST",
		Address:                "123 Testing Street",
		Applicant:              "Unit Test Inc",
		Category:               "Residential - Single Family",
		Description:            "NEW:RESIDENTIAL DEVELOPMENT (1 BUILDING), ACCESSORY RESIDENTIAL BUILDING (GARAGE) ",
		ProposedUseCode:        "T0001; T0002",
		PermittedDiscretion:    "Permitted",
		LandUseDistrict:        "R-C2",
		LandUseDistrictDesc:    "Residential - Contextual One/Two Dwelling",
		StatusCurrent:          "Under Review",
		AppliedDate:            "2023-04-08T00:00:00.000",
		CommunityCode:          "KIL",
		CommunityName:          "KILLARNEY/GLENGARRY",
		Ward:                   "8",
		Quadrant:               "SW",
		Latitude:               "51.034",
		Longitude:              "-114.122",
		LocationCount:          "1",
		LocationTypes:          "Titled Parcel",
		LocationAddresses:      "123 Testing Street",
		LocationsGeoJSON:       "{\"type\":\"Point\",\"coordinates\":[-114.1224605,51.0338002]}",
		LocationsWKT:           "POINT (-114.12246048560655 51.03380020179407)",
		DecisionDate:           "",
		MustCommenceDate:       "",
		Decision:               "",
		DecisionBy:             "",
		ReleaseDate:            "",
		GithubDiscussionClosed: true,
	}
	fetchedPermit := DevelopmentPermit{
		PermitNum:              "DP2023-TEST",
		Address:                "123 Testing Street",
		Applicant:              "Unit Test Inc",
		Category:               "Residential - Single Family",
		Description:            "NEW:RESIDENTIAL DEVELOPMENT (1 BUILDING), ACCESSORY RESIDENTIAL BUILDING (GARAGE) ",
		ProposedUseCode:        "T0001; T0002",
		PermittedDiscretion:    "Permitted",
		LandUseDistrict:        "R-C2",
		LandUseDistrictDesc:    "Residential - Contextual One/Two Dwelling",
		StatusCurrent:          "Under Review",
		AppliedDate:            "2023-04-08T00:00:00.000",
		CommunityCode:          "KIL",
		CommunityName:          "KILLARNEY/GLENGARRY",
		Ward:                   "8",
		Quadrant:               "SW",
		Latitude:               "51.034",
		Longitude:              "-114.122",
		LocationCount:          "1",
		LocationTypes:          "Titled Parcel",
		LocationAddresses:      "123 Testing Street",
		LocationsGeoJSON:       "{\"type\":\"Point\",\"coordinates\":[-114.1224605,51.0338002]}",
		LocationsWKT:           "POINT (-114.12246048560655 51.03380020179407)",
		DecisionDate:           "",
		MustCommenceDate:       "",
		Decision:               "",
		DecisionBy:             "",
		ReleaseDate:            "",
		GithubDiscussionClosed: false,
	}

	fetchedPermits := []DevelopmentPermit{fetchedPermit}
	storedPermits := []DevelopmentPermit{storedPermit}

	actions := GetDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 1, len(actions))
	assert.Equal(t, "DP2023-TEST", actions[0].PermitNum)
	assert.Equal(t, "SKIP", actions[0].Action)
	assert.Equal(t, "", actions[0].Message)
}

func TestContains_DevelopmentPermitAction_CreateAction(t *testing.T) {
	fetchedPermit := DevelopmentPermit{
		PermitNum:              "DP2023-TEST",
		Address:                "123 Testing Street",
		Applicant:              "Unit Test Inc",
		Category:               "Residential - Single Family",
		Description:            "NEW:RESIDENTIAL DEVELOPMENT (1 BUILDING), ACCESSORY RESIDENTIAL BUILDING (GARAGE) ",
		ProposedUseCode:        "T0001; T0002",
		PermittedDiscretion:    "Permitted",
		LandUseDistrict:        "R-C2",
		LandUseDistrictDesc:    "Residential - Contextual One/Two Dwelling",
		StatusCurrent:          "Under Review",
		AppliedDate:            "2023-04-08T00:00:00.000",
		CommunityCode:          "KIL",
		CommunityName:          "KILLARNEY/GLENGARRY",
		Ward:                   "8",
		Quadrant:               "SW",
		Latitude:               "51.034",
		Longitude:              "-114.122",
		LocationCount:          "1",
		LocationTypes:          "Titled Parcel",
		LocationAddresses:      "123 Testing Street",
		LocationsGeoJSON:       "{\"type\":\"Point\",\"coordinates\":[-114.1224605,51.0338002]}",
		LocationsWKT:           "POINT (-114.12246048560655 51.03380020179407)",
		DecisionDate:           "",
		MustCommenceDate:       "",
		Decision:               "",
		DecisionBy:             "",
		ReleaseDate:            "",
		GithubDiscussionClosed: false,
	}

	fetchedPermits := []DevelopmentPermit{fetchedPermit}
	storedPermits := []DevelopmentPermit{}

	actions := GetDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 1, len(actions))
	assert.Equal(t, "DP2023-TEST", actions[0].PermitNum)
	assert.Equal(t, "CREATE", actions[0].Action)
	assert.Contains(t, actions[0].Message, "Development Permit Application")
}
