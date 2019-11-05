package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/docopt/docopt-go"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	// EXITSUCCESS represents a successful command line program status code.
	EXITSUCCESS = 0
	// EXITFAILURE represents the status code for a command line program that failled.
	EXITFAILURE = 1
	// USAGE is the documentation for the command line.
	USAGE = `Github reporter

Usage:
  gh-reporter issues (--owner=<owner> --repo=<repo> --since=<since> --to=<to>) [--state=<state>]
  gh-reporter cards (--owner=<owner> --repo=<repo> --column-id=<column_id>)
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

var version = "dev"

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
	fmt.Println(arguments)

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

	} else if arguments["cards"] == true {
		owner := arguments["--owner"].(string)
		repo := arguments["--repo"].(string)
		columnID, err := strconv.Atoi(arguments["--column-id"].(string))
		if err != nil {
			fmt.Printf("An error occured while converting column-id: %v\n", err)
		}

		err = reportCards(client, owner, repo, int64(columnID))
		if err != nil {
			fmt.Printf("An error occured while retrieving cards in project column: %v\n", err)
			os.Exit(EXITFAILURE)
		}
		os.Exit(EXITFAILURE)
	} else if arguments["--version"] == true {
		fmt.Println("version: ", version)
	}
	os.Exit(EXITSUCCESS)
}
