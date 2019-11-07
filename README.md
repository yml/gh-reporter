gh-reporter
=============

gh-reporter is command line interface to Github built on top of go-github

export your GITHUB api 

```
export GITHUB_TOKEN="<API TOKEN>
go get github.com/yml/gh-reporter
go install github.com/yml/gh-reporter
./gh-reporter -org=lincolnloop -since="2013-07-29T00:00:00Z" -to="2013-08-29T00:00:00Z"
```
