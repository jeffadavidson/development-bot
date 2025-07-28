
## About

A GO application that can run on a schedule and generate RSS feeds for new, updated and closed Development Permits and Land Use Designations for the community of Killarney in Calgary Alberta

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
- `development-permits.xml` - Generated RSS feed for development permits
- `rezoning-applications.xml` - Generated RSS feed for rezoning applications

The JSON files track what's been processed to detect new/updated permits. The XML files are the RSS feeds that can be served to users.

## RSS Feed URLs

When hosted on GitHub Pages, the RSS feeds will be available at:
- Development Permits: `https://username.github.io/repo-name/data/development-permits.xml`
- Rezoning Applications: `https://username.github.io/repo-name/data/rezoning-applications.xml`

## Future Plans

- Add web server to serve RSS feeds directly
- Add more detailed permit information to RSS descriptions
- Add geographic filtering options