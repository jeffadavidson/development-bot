package examinedata

import (
	"testing"

	"github.com/jeffadavidson/development-bot/objects/developmentpermit"
	"github.com/stretchr/testify/assert"
)

func TestContains_DevelopmentPermitAction_NoUpdate(t *testing.T) {
	storedPermit := developmentpermit.DevelopmentPermit{
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
	fetchedPermit := developmentpermit.DevelopmentPermit{
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

	fetchedPermits := []developmentpermit.DevelopmentPermit{fetchedPermit}
	storedPermits := []developmentpermit.DevelopmentPermit{storedPermit}

	actions := getDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 0, len(actions))
}

func TestContains_DevelopmentPermitAction_StatusUpdate(t *testing.T) {
	storedPermit := developmentpermit.DevelopmentPermit{
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
	fetchedPermit := developmentpermit.DevelopmentPermit{
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

	fetchedPermits := []developmentpermit.DevelopmentPermit{fetchedPermit}
	storedPermits := []developmentpermit.DevelopmentPermit{storedPermit}

	actions := getDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 1, len(actions))
	assert.Equal(t, "DP2023-TEST", actions[0].PermitNum)
	assert.Contains(t, actions[0].Action, "UPDATE")
	assert.Contains(t, actions[0].Message, "Status")
}

func TestContains_DevelopmentPermitAction_DecisionUpdate(t *testing.T) {
	storedPermit := developmentpermit.DevelopmentPermit{
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
	fetchedPermit := developmentpermit.DevelopmentPermit{
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

	fetchedPermits := []developmentpermit.DevelopmentPermit{fetchedPermit}
	storedPermits := []developmentpermit.DevelopmentPermit{storedPermit}

	actions := getDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 1, len(actions))
	assert.Equal(t, "DP2023-TEST", actions[0].PermitNum)
	assert.Contains(t, actions[0].Action, "UPDATE")
	assert.Contains(t, actions[0].Message, "Decision")
}

func TestContains_DevelopmentPermitAction_MustCommenceByUpdate(t *testing.T) {
	storedPermit := developmentpermit.DevelopmentPermit{
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
	fetchedPermit := developmentpermit.DevelopmentPermit{
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

	fetchedPermits := []developmentpermit.DevelopmentPermit{fetchedPermit}
	storedPermits := []developmentpermit.DevelopmentPermit{storedPermit}

	actions := getDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 1, len(actions))
	assert.Equal(t, "DP2023-TEST", actions[0].PermitNum)
	assert.Contains(t, actions[0].Action, "UPDATE")
	assert.Contains(t, actions[0].Message, "Must Commence")
}

func TestContains_DevelopmentPermitAction_MultipleUpdates(t *testing.T) {
	storedPermit := developmentpermit.DevelopmentPermit{
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
	fetchedPermit := developmentpermit.DevelopmentPermit{
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

	fetchedPermits := []developmentpermit.DevelopmentPermit{fetchedPermit}
	storedPermits := []developmentpermit.DevelopmentPermit{storedPermit}

	actions := getDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 1, len(actions))
	assert.Equal(t, "DP2023-TEST", actions[0].PermitNum)
	assert.Contains(t, actions[0].Action, "UPDATE")
	assert.Contains(t, actions[0].Message, "Status")
	assert.Contains(t, actions[0].Message, "Decision")
	assert.Contains(t, actions[0].Message, "Must Commence")
}

func TestContains_DevelopmentPermitAction_StatusUpdateReleased(t *testing.T) {
	storedPermit := developmentpermit.DevelopmentPermit{
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
	fetchedPermit := developmentpermit.DevelopmentPermit{
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

	fetchedPermits := []developmentpermit.DevelopmentPermit{fetchedPermit}
	storedPermits := []developmentpermit.DevelopmentPermit{storedPermit}

	actions := getDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 1, len(actions))
	assert.Equal(t, "DP2023-TEST", actions[0].PermitNum)
	assert.Contains(t, actions[0].Action, "CLOSE")
	assert.Contains(t, actions[0].Message, "Status")
	assert.Contains(t, actions[0].Message, "Closing")
}

func TestContains_DevelopmentPermitAction_DiscussionClosed(t *testing.T) {
	storedPermit := developmentpermit.DevelopmentPermit{
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
	fetchedPermit := developmentpermit.DevelopmentPermit{
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

	fetchedPermits := []developmentpermit.DevelopmentPermit{fetchedPermit}
	storedPermits := []developmentpermit.DevelopmentPermit{storedPermit}

	actions := getDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 1, len(actions))
	assert.Equal(t, "DP2023-TEST", actions[0].PermitNum)
	assert.Equal(t, "SKIP", actions[0].Action)
	assert.Equal(t, "", actions[0].Message)
}

func TestContains_DevelopmentPermitAction_CreateAction(t *testing.T) {
	fetchedPermit := developmentpermit.DevelopmentPermit{
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

	fetchedPermits := []developmentpermit.DevelopmentPermit{fetchedPermit}
	storedPermits := []developmentpermit.DevelopmentPermit{}

	actions := getDevelopmentPermitActions(fetchedPermits, storedPermits)
	assert.Equal(t, 1, len(actions))
	assert.Equal(t, "DP2023-TEST", actions[0].PermitNum)
	assert.Equal(t, "CREATE", actions[0].Action)
	assert.Contains(t, actions[0].Message, "Development Permit Application")
}
