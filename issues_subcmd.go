package main

import (
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"github.com/google/go-github/github"
	"os"
	"time"
)

var cmdIssues = &Command{
	UsageLine: "issues [-state open] [-since] [-to]",
	Short:     "List gihub issues",
	Long: `
List gihub issues

#TODO add more here later
`,
}

var Org = cmdIssues.Flag.String("org", "", "Github Organisation you want to query against (ie lincolnloop)")
var Since = cmdIssues.Flag.String("since", "", "Since date (ie 2013-07-29T00:00:00Z)")
var To = cmdIssues.Flag.String("to", "", "To date (ie 2013-08-09T00:00:00Z)")
var State = cmdIssues.Flag.String("state", "", "State  open|close|all")

func init() {
	cmdIssues.Run = runIssues
}

type Page struct {
	Number  int
	Next    int
	Last    int
	Fetched bool
	Result  []github.Issue
}

type Pager struct {
	Pages []*Page
}

func (pager *Pager) Add(page *Page) {
	pager.Pages = append(pager.Pages, page)
}

func (pager *Pager) IsFetched() bool {
	for _, page := range pager.Pages {
		if page.Fetched {
			continue
		} else {
			return false
		}
	}
	return true
}

func fetchPageIssue(client *github.Client, opts github.IssueListOptions, page *Page) (err error) {
	if err != nil {
		return err
	}

	opts.Page = page.Number

	issues, response, err := client.Issues.ListByOrg(*Org, &opts)
	if err == nil {
		page.Fetched = true
		page.Result = issues
	}
	page.Next = response.NextPage
	page.Last = response.LastPage
	return err
}

func IssuePager(opts github.IssueListOptions) (pager Pager, err error) {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	}
	client := github.NewClient(t.Client())
	pager = Pager{}
	page := &Page{Number: 1}
	pager.Add(page)
	err = fetchPageIssue(client, opts, page)
	if err != nil {
		return pager, err
	}

	for i := page.Next; i <= page.Last; i++ {
		page := &Page{Number: i}
		pager.Add(page)
		go func(page *Page) {
			fetchPageIssue(client, opts, page)
		}(page)
	}
	// Wait until all the Pages are fetched
	for !pager.IsFetched() {
		time.Sleep(1 * time.Second)
	}
	return pager, nil
}

func StringifyIssue(issue github.Issue) string {
	return fmt.Sprintf("%s #%d %s %s %s %s\n",
		issue.UpdatedAt.Format(time.RFC822), issue.Number,
		issue.State, issue.User.Login, issue.Labels, issue.Title)
}

func runIssues(cmd *Command, args []string) {
	since, err := time.Parse(time.RFC3339, *Since)
	if err != nil {
		fmt.Println("An error occured while parsing the `since` date:", err)
		cmd.Usage()
		setExitStatus(EXIT_FAILURE)
	}
	to, err := time.Parse(time.RFC3339, *To)
	if err != nil {
		fmt.Println("An error occured while parsing the `To` date:", err)
		cmd.Usage()
		setExitStatus(EXIT_FAILURE)
	}

	// Fetch the issues for the period
	opts := github.IssueListOptions{
		Sort:      "updated",
		Direction: "desc",
		Filter:    "all",
		State:     *State,
		Since:     since,
	}

	pager, err := IssuePager(opts)
	if err != nil {
		fmt.Println("An error occured while querying github: ", err)
		setExitStatus(EXIT_FAILURE)
	}

	issueCount := 0

	for _, page := range pager.Pages {
		for _, issue := range page.Result {
			if to.After(*issue.UpdatedAt) {
				issueCount++
				fmt.Printf(StringifyIssue(issue))
			}
		}
	}
}
