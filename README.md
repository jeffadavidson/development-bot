
## About

A GO application that can run on a schedule and generate RSS feeds for new, updated and closed Development Permits and Land Use Designations for the community of Killarney in Calgary Alberta. The application tracks activity from the last 3 months to focus on current and relevant development activity.

We care about Development Permits and Land Use Redesignations in Killarney/Glengarry as well as on the edges of our boarders. The bounding box we care about is defined as:

* Northwest Corner
  * Latitude:  51.038912
  * Longitude: -114.142638
* Northeast Corner
  * Latitude:  51.038912
  * Longitude: -114.117927
* Southwest Corner
  * Latitude:  51.022361
  * Longitude: -114.142638
* Southeast Corner
  * Latitude:  51.022361
  * Longitude: -114.117927

## Execution

### Running on local machine:

Assumes Go is installed

```Shell
go test -v ./...
go run main.go
```

The application will:
- Fetch the latest data from Calgary Open Data
- Compare with stored data in `./data/` 
- Create/update RSS XML files in `./data/development-permits.xml` and `./data/rezoning-applications.xml`
- Display console output about what RSS entries were created/updated
## Persistent storage

The bot loads files from persistent storage in the `data/` directory:
- `development-permits.json` - Stores processed development permit data
- `rezoning-applications.json` - Stores processed rezoning application data

The combined RSS feed is generated at:
- `killarney-development.xml` - Combined RSS feed for all development activity (in root directory)

The JSON files track what's been processed to detect new/updated permits and store:
- **Stable RSS GUIDs** for each permit/application 
- **Complete state history** tracking all status changes over time with timestamps

The XML file is the combined RSS feed with both development permits (🏗️) and rezoning applications (🏛️) in chronological order. Each permit/application has a stable GUID that persists through status changes, ensuring RSS readers see updates rather than duplicate entries.

### State History
Each permit/application maintains a complete audit trail of status changes:
```json
"state_history": [
  {
    "status": "submitted",
    "timestamp": "2025-07-28T09:41:04-06:00"
  },
  {
    "status": "under review", 
    "timestamp": "2025-07-29T10:15:32-06:00",
    "decision": "approved"
  }
]
```

This enables full lifecycle reconstruction and detailed reporting on permit progression.

## RSS Feed URL

When hosted on GitHub Pages, the RSS feed will be available at:
- **All Development Activity**: `https://kgca-development.github.io/development-bot/killarney-development.xml`
- **Web Interface**: `https://kgca-development.github.io/development-bot/`

Each RSS entry includes comprehensive metadata and links directly to Calgary's [Development Map](https://developmentmap.calgary.ca/) for that specific permit or application.

### Enhanced RSS Metadata
Each RSS item now includes:
- **Title**: Shows current status (e.g., "🏗️ Development Permit (Under Review): DP2025-12345 - 123 Main St")
- **Description**: Full permit details in markdown format
- **Link**: Direct link to Calgary Development Map
- **Category**: "Development Permit" or "Land Use Rezoning"
- **Author**: Applicant name (when available)
- **Source**: "City of Calgary Open Data"
- **Comments**: Link to Development Map comments section
- **Publication Date**: Application submission date
- **GUID**: Stable unique identifier that persists through status changes

**Status-Based Titles**: Titles automatically update to reflect the current permit status (Hold, Under Review, In Circulation, Released, Cancelled, etc.) rather than static "New" or "Closed" labels.

## Future Plans
- Add web server to serve RSS feeds directly
- Add geographic filtering options
- Enhanced filtering by permit type or status