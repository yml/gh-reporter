package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

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
  gh-reporter issues [(--url=<URL>)|(--owner=<owner> --repo=<repo> )] [--since=<since> --to=<to> --state=<state>]
  gh-reporter cards [(--url=<URL>)|(--owner=<owner> --repo=<repo> --column-id=<column_id>)] [--title]
  gh-reporter -h | --help
  gh-reporter --version

Options:
  -h --help  # Show this screen.
  --version  # Show version.
  --url <URL>  # Github URL (ie https://github.com/yml/gh-reporter/issues)
  --owner <onwer>  # Github owner you want to query against (ie yml or lincolnloop)
  --repo <repo>  # Github repo you want to query against
  --since <since>  # Since date (ie 2019-07-29T00:00:00Z) [default: ]
  --to <to>  # To date (ie 2019-10-29T00:00:00Z) [default: ]
  --state <state>  # State  open|closed|all [default: all]
  --project-id <project>  # Project id
  --column-id <column_id>  # Column id
  --title  # Print out the ticket title
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

func exitWithError(msgFormat string, err error) {
	fmt.Printf(msgFormat, err)
	os.Exit(EXITFAILURE)
}

func main() {
	arguments, _ := docopt.ParseDoc(USAGE)
	fmt.Println(arguments)

	accessToken := os.Getenv("GITHUB_TOKEN")
	client := NewGithubClient(accessToken)
	var (
		withTitle bool
		owner     string
		repo      string
		err       error
		ghURL     *url.URL
	)

	if arguments["--title"] == true {
		withTitle = true
	} else {
		withTitle = false
	}
	fmt.Println(withTitle)

	if arguments["issues"] == true {
		var (
			since string
			to    string
			state string
		)
		if arguments["--url"] != nil {
			ghURL, err = url.Parse(arguments["--url"].(string))
			if err != nil {
				exitWithError("An error occured while parsing the URL: %v\n", err)
			}
			if ghURL.Hostname() != "github.com" {
				exitWithError("An error occured while parsing the URL: %v\n", fmt.Errorf("URL must be on the github.com domain and not: %s", ghURL.Hostname()))
			}

			fmt.Sscanf(strings.ReplaceAll(ghURL.Path, "/", " "), "%s %s issues", &owner, &repo)
		} else {
			owner = arguments["--owner"].(string)
			repo = arguments["--repo"].(string)
			since = arguments["--since"].(string)
			to = arguments["--to"].(string)
			state = arguments["--state"].(string)
		}

		err := reportIssues(client, owner, repo, since, to, state)
		if err != nil {
			exitWithError("An error occured while retrieving github issues: %v\n", err)
		}

	} else if arguments["cards"] == true {
		var (
			columnID int
		)
		if arguments["--url"] != nil {
			ghURL, err = url.Parse(arguments["--url"].(string))
			if err != nil {
				exitWithError("An error occured while parsing the URL: %v\n", err)
			}
			if ghURL.Hostname() != "github.com" {
				exitWithError("An error occured while parsing the URL: %v\n", fmt.Errorf("URL must be on the github.com domain and not: %s", ghURL.Hostname()))
			}

			fmt.Sscanf(strings.ReplaceAll(ghURL.Path, "/", " "), "%s %s projects %s", &owner, &repo)
			fmt.Sscanf(ghURL.Fragment, "column-%d", &columnID)

		} else {
			owner = arguments["--owner"].(string)
			repo = arguments["--repo"].(string)
			columnID, err = strconv.Atoi(arguments["--column-id"].(string))
			if err != nil {
				exitWithError("An error occured while converting column-id: %v\n", err)
			}
		}
		fmt.Println(owner, repo, columnID)
		err = reportCards(client, owner, repo, int64(columnID), withTitle)
		if err != nil {
			exitWithError("An error occured while retrieving cards in project column: %v\n", err)
		}
	} else if arguments["--version"] == true {
		fmt.Println("version: ", version)
	}
	os.Exit(EXITSUCCESS)
}
