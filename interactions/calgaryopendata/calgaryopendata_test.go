package calgaryopendata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
