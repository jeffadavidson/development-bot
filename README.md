
## About

A GO application that can run on a schedule and update a private GitHub discussion board with new, updated and closed Development Permits and Land Use Designations for the community of Killarney in Calgary Alberta



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

### Setting Environment Variable

| Variable Name                  | Purpose                                                             | Allowed Values              |
| ------------------------------ | ------------------------------------------------------------------- | --------------------------- |
| `DEVELOPMENT_BOT_GITHUB_PAT`   | The development bot GitHub PAT with Repo and Discussion Permissions | A Valid GitHub Classic PAT  |
| `DEVELOPMENT_BOT_RUNMODE`      | The development bot run mode. Will actions actually be performed    | `PRODUCTION`, `DEVELOPMENT` |

### Running on local machine:

Assumes Go is installed

Unix:

```Shell
export DEVELOPMENT_BOT_GITHUB_PAT=<somekey>
export DEVELOPMENT_BOT_RUNMODE=DEVELOPMENT
printenv DEVELOPMENT_BOT_GITHUB_PAT
printenv DEVELOPMENT_BOT_RUNMODE

go test -v ./...
go run main.go
```

Windows:

```Shell
$env:DEVELOPMENT_BOT_GITHUB_PAT="<token>"
$env:DEVELOPMENT_BOT_RUNMODE="DEVELOPMENT"
$env:DEVELOPMENT_BOT_GITHUB_PAT
$env:DEVELOPMENT_BOT_RUNMODE
go test -v ./...
go run main.go
```
## Persistent storage

The bot loads a files from persistent storage, the directory `data/`. This should/could be cloud storage but I am cheep so I keep the data stored in the repository itself, its not a lot so its passable, but means all hell will break lose if multiple instances run at the same time.

## Future Plans

* Add labels to discussions
* Trigger issue create and assign to Jeff for LOCs in <some> state
* Archive/Unarchive on Hold permits
* v2: Community bots
  * Water main breaks
  * Traffic incidents
  * development permits
  * Property assessments
  * construction detours
  * road construction
  * crime
  * trees?
  * weather and Rain