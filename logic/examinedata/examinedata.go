package examinedata

import (
	"fmt"

	"github.com/jeffadavidson/development-bot/interactions/rssfeed"
	"github.com/jeffadavidson/development-bot/objects/developmentpermit"
	"github.com/jeffadavidson/development-bot/objects/rezoningapplications"
)

func ManualInit() error {
	// No initialization needed for RSS feeds
	return nil
}

// ProcessAllDevelopmentActivity - Evaluates both development permits and rezoning applications and generates a combined RSS feed
func ProcessAllDevelopmentActivity() error {
	// Load or create combined RSS feed
	rss, err := rssfeed.GetOrCreateRSSFeed(
		"./killarney-development.xml",
		"Killarney Development Activity",
		"All development permits and land use rezoning applications for the Killarney neighborhood in Calgary",
		"https://calgary.ca/development",
	)
	if err != nil {
		return fmt.Errorf("failed to load RSS feed: %v", err)
	}

	// Process development permits
	dpActions, dpErr := developmentpermit.EvaluateDevelopmentPermits(rss)
	if dpErr != nil {
		return fmt.Errorf("failed to process development permits: %v", dpErr)
	}

	// Process rezoning applications  
	raActions, raErr := rezoningapplications.EvaluateRezoningApplications(rss)
	if raErr != nil {
		return fmt.Errorf("failed to process rezoning applications: %v", raErr)
	}

	// Trim RSS feed to keep only recent items (increased since we have both types)
	rss.TrimToMaxItems(200)

	// Save combined RSS feed
	if err := rssfeed.SaveRSSFeed(rss, "./killarney-development.xml"); err != nil {
		return fmt.Errorf("failed to save RSS feed: %v", err)
	}

	fmt.Printf("Combined RSS feed processed with %d development permit actions and %d rezoning application actions\n", 
		len(dpActions), len(raActions))

	return nil
}
