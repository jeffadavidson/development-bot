package examinedata

import (
	"github.com/jeffadavidson/development-bot/objects/developmentpermit"
	"github.com/jeffadavidson/development-bot/objects/rezoningapplications"
)

func ManualInit() error {
	// No initialization needed for RSS feeds
	return nil
}

// DevelopmentPermits - Evaluates Development permits and generates RSS feed
func DevelopmentPermits() error {
	return developmentpermit.EvaluateDevelopmentPermits()
}

// RezoningApplication - Evaluates Rezoning Applications and generates RSS feed
func RezoningApplication() error {
	return rezoningapplications.EvaluateRezoningApplications()
}
