package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	EXITSUCCESS = 0
	EXITFAILURE = 1
	USAGE       = `Github reporter

Usage:
  gh-reporter issues (--owner=<owner> --repo=<repo> --since=<since> --to=<to>) [--state=<state>]
  gh-reporter project (--owner=<owner> --repo=<repo> --project-id=<project_id> --column-id=<column_id>)
  gh-reporter -h | --help
  gh-reporter --version

Options:
  -h --help  # Show this screen.
  --version  # Show version.
  --owner <onwer>  # Github owner you want to query against (ie yml or lincolnloop)
  --repo <repo>  # Github repo you want to query against
  --since <since>  # Since date (ie 2019-07-29T00:00:00Z)
  --to <to>  # To date (ie 2019-10-29T00:00:00Z)
  --state <state>  # State  open|closed|all [default: all]
  --project-id <project>  # Project id
  --column-id <column_id>  # Column id
`
)

//NewGithubClient return an initialize Github Client
func NewGithubClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func main() {
	arguments, _ := docopt.ParseDoc(USAGE)
	// fmt.Println(arguments)

	accessToken := os.Getenv("GITHUB_TOKEN")
	client := NewGithubClient(accessToken)
	if arguments["issues"] == true {
		owner := arguments["--owner"].(string)
		repo := arguments["--repo"].(string)
		since := arguments["--since"].(string)
		to := arguments["--to"].(string)
		state := arguments["--state"].(string)

		err := runIssues(client, owner, repo, since, to, state)
		if err != nil {
			fmt.Printf("An error occured while retrieving github issues: %v\n", err)
			os.Exit(EXITFAILURE)
		}

	} else if arguments["projects"] == true {
		fmt.Println("Project snapshot is Not Implmented")
		os.Exit(EXITFAILURE)
	}
	os.Exit(EXITSUCCESS)
}
