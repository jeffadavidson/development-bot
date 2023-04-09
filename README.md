
## About

## Execution

### Environment Variables

development-bot uses environment variables to provide secret configuration as well as to control execution flow.

* `DEVELOPMENT_BOT_GITHUB_PAT` - A classic personal access token for accessing Github. Need Repo and Discussion(read:write) permissions

```unix
export DEVELOPMENT_BOT_GITHUB_PAT=<somekey>
printenv DEVELOPMENT_BOT_GITHUB_PAT

go run main.go
```


```windows
$env:DEVELOPMENT_BOT_GITHUB_PAT="<token>"
$env:DEVELOPMENT_BOT_GITHUB_PAT


go run main.go

go build main.go

go test -v ./... 
```

## Configuration GitHub Actions

For the code to execute in GitHub Actions we need to add a secret to that repos GitHub actions

* Get PAT
  * Login as the owner of the repo you want to own/create discussions in. Likely `kgca-development/killarney-development`
  * Navigate to personal access tokens (Settings > Developer Settings > Personal Access Token (classic))
  * Generate a new Classic Personal Access Token, hold onto that in a scratch pat
* Set PAT
  * Navigate to the repo you want to run the development-bot in
  * Navigate Settings > Actions > Secrets
  * Add a secret called `DEVELOPMENT_BOT_GITHUB_PAT`