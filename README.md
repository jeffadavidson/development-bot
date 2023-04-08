
## About

## Execution

### Environment Variables

development-bot uses environment variables to provide secret configuration as well as to control execution flow.

* `GITHUB_PAT` - A classic personal access token for accessing Github. Need Repo and Discussion(read:write) permissions

```unix
export GITHUB_PAT=<somekey>
printenv GITHUB_PAT

go run main.go
```


```windows
$env:GITHUB_PAT="<token>"
$env:GITHUB_PAT


go run main.go

go build main.go

go test -v ./... 
```