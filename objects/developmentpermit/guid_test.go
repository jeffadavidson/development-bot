package developmentpermit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateUniqueGUID_Consistency(t *testing.T) {
	permitNum := "DP2025-12345"
	permitType := "development-permit"
	
	// Generate GUID multiple times
	guid1 := generateUniqueGUID(permitNum, permitType)
	guid2 := generateUniqueGUID(permitNum, permitType)
	guid3 := generateUniqueGUID(permitNum, permitType)
	
	// Should be consistent
	assert.Equal(t, guid1, guid2)
	assert.Equal(t, guid2, guid3)
	
	// Should be 32 characters (16 bytes hex)
	assert.Len(t, guid1, 32)
}

func TestGenerateUniqueGUID_Different(t *testing.T) {
	// Different permit numbers should generate different GUIDs
	guid1 := generateUniqueGUID("DP2025-12345", "development-permit")
	guid2 := generateUniqueGUID("DP2025-54321", "development-permit")
	
	assert.NotEqual(t, guid1, guid2)
	
	// Different permit types should generate different GUIDs
	guid3 := generateUniqueGUID("DP2025-12345", "rezoning-application")
	
	assert.NotEqual(t, guid1, guid3)
}

func TestUpdateStateHistory_NewPermit(t *testing.T) {
	permit := &DevelopmentPermit{
		PermitNum:     "DP2025-12345",
		StatusCurrent: "Under Review",
		Decision:      stringPointer("Pending"),
	}
	
	updateStateHistory(permit, nil)
	
	assert.Len(t, permit.StateHistory, 1)
	assert.Equal(t, "under review", permit.StateHistory[0].Status)
	assert.Equal(t, "Pending", permit.StateHistory[0].Decision)
	assert.NotEmpty(t, permit.StateHistory[0].Timestamp)
}

func TestUpdateStateHistory_StatusChange(t *testing.T) {
	// Create stored permit with existing history
	storedPermit := &DevelopmentPermit{
		PermitNum:     "DP2025-12345",
		StatusCurrent: "Under Review",
		StateHistory: []StateChange{
			{
				Status:    "under review",
				Timestamp: "2025-07-01T12:00:00Z",
				Decision:  "Pending",
			},
		},
	}
	
	// Create fetched permit with new status
	fetchedPermit := &DevelopmentPermit{
		PermitNum:     "DP2025-12345",
		StatusCurrent: "Approved",
		Decision:      stringPointer("Approved"),
	}
	
	updateStateHistory(fetchedPermit, storedPermit)
	
	// Should have 2 entries
	assert.Len(t, fetchedPermit.StateHistory, 2)
	
	// First entry should be preserved
	assert.Equal(t, "under review", fetchedPermit.StateHistory[0].Status)
	assert.Equal(t, "2025-07-01T12:00:00Z", fetchedPermit.StateHistory[0].Timestamp)
	
	// Second entry should be new
	assert.Equal(t, "approved", fetchedPermit.StateHistory[1].Status)
	assert.Equal(t, "Approved", fetchedPermit.StateHistory[1].Decision)
	assert.NotEmpty(t, fetchedPermit.StateHistory[1].Timestamp)
}

func TestUpdateStateHistory_NoChange(t *testing.T) {
	// Create stored permit 
	storedPermit := &DevelopmentPermit{
		PermitNum:     "DP2025-12345",
		StatusCurrent: "Under Review",
		StateHistory: []StateChange{
			{
				Status:    "under review",
				Timestamp: "2025-07-01T12:00:00Z",
			},
		},
	}
	
	// Create fetched permit with same status
	fetchedPermit := &DevelopmentPermit{
		PermitNum:     "DP2025-12345",
		StatusCurrent: "Under Review",
	}
	
	updateStateHistory(fetchedPermit, storedPermit)
	
	// Should still have 1 entry (no change)
	assert.Len(t, fetchedPermit.StateHistory, 1)
	assert.Equal(t, "under review", fetchedPermit.StateHistory[0].Status)
}

func TestGetStateHistorySummary(t *testing.T) {
	permit := DevelopmentPermit{
		PermitNum: "DP2025-12345",
		StateHistory: []StateChange{
			{
				Status:    "submitted",
				Timestamp: "2025-07-01T12:00:00Z",
				Decision:  "",
			},
			{
				Status:    "under review",
				Timestamp: "2025-07-15T10:30:00Z",
				Decision:  "Pending",
			},
			{
				Status:    "approved",
				Timestamp: "2025-07-20T14:15:00Z",
				Decision:  "Approved",
			},
		},
	}
	
	summary := permit.GetStateHistorySummary()
	
	assert.Contains(t, summary, "DP2025-12345 lifecycle:")
	assert.Contains(t, summary, "1. Submitted")
	assert.Contains(t, summary, "2. Under Review")
	assert.Contains(t, summary, "3. Approved")
	assert.Contains(t, summary, "(Decision: Pending)")
	assert.Contains(t, summary, "(Decision: Approved)")
}

func TestGetStateHistorySummary_Empty(t *testing.T) {
	permit := DevelopmentPermit{
		PermitNum:    "DP2025-12345",
		StateHistory: []StateChange{},
	}
	
	summary := permit.GetStateHistorySummary()
	
	assert.Equal(t, "No state history available", summary)
}

// Helper function for tests
func stringPointer(s string) *string {
	return &s
} 