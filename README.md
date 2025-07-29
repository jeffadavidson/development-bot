
# Development Bot 🏗️

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org/)
[![Tests](https://github.com/kgca-development/development-bot/workflows/Tests/badge.svg)](https://github.com/kgca-development/development-bot/actions)

An automated civic technology tool that generates RSS feeds for local development activity. Originally built for the Killarney/Glengarry community in Calgary, but designed to be adaptable for any city with open data APIs.

## 🎯 What It Does

Tracks and publishes development permits and land use rezoning applications through:
- **📡 RSS Feed**: Automatically updated feed of development activity
- **🏗️ Development Permits**: Building permits, renovations, new construction
- **🏛️ Rezoning Applications**: Land use changes and redesignations  
- **📊 State History**: Complete audit trail of permit status changes
- **🔄 Daily Updates**: Runs automatically at 6AM MT via GitHub Actions

The application tracks activity from the last 3 months to focus on current and relevant development activity.

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

## 🚀 How to Run

### Prerequisites
- Go 1.24+ installed
- Internet connection (for Calgary Open Data API access)

### Local Development
```bash
# Clone the repository
git clone https://github.com/kgca-development/development-bot.git
cd development-bot

# Install dependencies
go mod download

# Run the application
go run main.go
```

### What happens when you run it:
1. **Fetches data** from Calgary Open Data API for development permits and rezoning applications
2. **Compares** with stored data in `./data/` directory
3. **Generates RSS feed** at `./output/killarney-development.xml`
4. **Updates stored data** for future comparisons
5. **Console output** shows what entries were created/updated:
   ```
   Development Permit DP2025-12345:
       Creating RSS feed entry...
       Created RSS feed entry!
   Combined RSS feed processed with 5 development permit actions and 2 rezoning application actions
   ```

### Console Output Meanings:
- **"Creating RSS feed entry"**: New permit/application found
- **"Updating RSS feed entry"**: Existing permit status changed
- **"0 actions"**: No changes detected (normal for subsequent runs)

## 🧪 How to Test

### Run All Tests
```bash
# Run all tests with verbose output
go test ./... -v

# Run tests for specific modules
go test ./interactions/rssfeed/ -v
go test ./objects/developmentpermit/ -v
go test ./objects/rezoningapplications/ -v
```

### Test Categories
- **RSS Feed Tests**: XML generation, namespace handling, item updates
- **Development Permit Tests**: Data parsing, action detection, timestamp handling
- **Rezoning Application Tests**: Status changes, close detection
- **Integration Tests**: End-to-end workflows

### Manual Testing
```bash
# Clear stored data to test full regeneration
rm -rf data/*.json output/*.xml

# Run and verify all entries are created
go run main.go

# Run again to verify no duplicate actions
go run main.go
```

### Validate RSS Output
```bash
# Check the generated RSS feed
cat output/killarney-development.xml

# Verify XML structure is valid
xmllint --noout output/killarney-development.xml
```

## 🚀 How Deployment Works

### GitHub Actions Workflow
The application automatically deploys via GitHub Actions (`.github/workflows/development-bot-runs.yml`):

#### Triggers:
- **Daily**: 12:00 PM UTC (6:00 AM MT) 
- **Manual**: Workflow dispatch from GitHub UI
- **Push**: When code is pushed to main branch

#### Deployment Process:
1. **Setup Environment**:
   ```yaml
   - Checkout code with deploy key
   - Install Go 1.24
   - Download dependencies
   ```

2. **Execute Bot**:
   ```bash
   go run main.go
   ```

3. **Commit Changes**:
   ```bash
   git add data/ output/
   git commit -m "Update development data and RSS feed"
   git push origin main
   ```

4. **Deploy to GitHub Pages**:
   ```bash
   # Copy RSS feed to Pages directory
   cp output/killarney-development.xml _site/
   # Deploy to https://kgca-development.github.io/development-bot/
   ```

### Deploy Key Authentication
- Uses `DEVELOPMENT_BOT_DEPLOY_KEY` secret for git operations
- Bypasses branch protection rules
- Enables automated commits to main branch

### Live RSS Feed URL
When deployed, the RSS feed is available at:
- **RSS Feed**: [`https://kgca-development.github.io/development-bot/killarney-development.xml`](https://kgca-development.github.io/development-bot/killarney-development.xml)
- **Web Interface**: [`https://kgca-development.github.io/development-bot/`](https://kgca-development.github.io/development-bot/)

### Monitoring Deployment
- Check [GitHub Actions](https://github.com/kgca-development/development-bot/actions) for run status
- View commit history for data updates
- Monitor RSS feed for new entries

## Persistent Storage

The bot stores data in two locations:

### `./data/` Directory (Version Controlled)
- `development-permits.json` - Processed development permit data with state history
- `rezoning-applications.json` - Processed rezoning application data with state history

### `./output/` Directory (Version Controlled)
- `killarney-development.xml` - Combined RSS feed for all development activity

### Data Structure
JSON files track what's been processed and store:
- **Stable RSS GUIDs** for each permit/application 
- **Complete state history** tracking all status changes over time with timestamps
- **Full permit data** for comparison on subsequent runs

### State History Example
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

## RSS Feed Features

### Enhanced RSS Metadata
Each RSS item includes:
- **Title**: Shows permit type and address (e.g., "🏗️ Development Permit: DP2025-12345 - 123 Main St")
- **Description**: Full permit details with rich HTML formatting
- **Link**: Direct link to Calgary Development Map
- **Category**: "Development Permit" or "Land Use Rezoning"
- **Author**: Applicant name (when available)
- **Source**: "City of Calgary Open Data"
- **Comments**: Link to Development Map comments section
- **Publication Date**: Most recent timestamp from permit data (not midnight UTC!)
- **GUID**: Stable unique identifier that persists through status changes

### XML Namespace Compliance
- Proper `xmlns:content` namespace declaration
- Valid `content:encoded` elements for rich content
- Validates in all RSS readers and browsers

## 🌍 Adapting for Your City

This bot can be adapted for any city with open data APIs:

1. **Fork this repository**
2. **Update configuration**:
   ```yaml
   # config.yaml
   neighborhood:
     name: "Your Neighborhood"
     bounding-box:
       north-latitude: 51.038912
       east-longitude: -114.117927
       south-latitude: 51.022361
       west-longitude: -114.142638
   ```
3. **Modify API endpoints** in `interactions/calgaryopendata/`
4. **Adjust data parsing** for your city's JSON structure
5. **Enable GitHub Actions** and **GitHub Pages** in your fork

## 🤝 Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

- 🐛 **Bug Reports**: Use GitHub Issues
- 💡 **Feature Requests**: Use GitHub Issues with "enhancement" label  
- 🔧 **Code Contributions**: Fork, branch, test, and submit PR
- 📖 **Documentation**: Help improve setup guides and examples

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🏙️ Community Impact

This tool helps communities:
- **Stay Informed**: Know what's being built nearby
- **Engage Actively**: Participate in local development discussions
- **Track Changes**: Monitor permit status from application to completion
- **Access Data**: Make government data more accessible

## 🔮 Future Plans
- Add web server to serve RSS feeds directly
- Add geographic filtering options
- Enhanced filtering by permit type or status
- Multi-city support in single deployment
- Email notifications for RSS updates