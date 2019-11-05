package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/github"
)

type Page struct {
	Number  int
	Next    int
	Last    int
	Fetched bool
	Result  []*github.Issue
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

// GhIssues hold the logic to fetch information related to a github project.
type GhIssues struct {
	Owner string
	Repo  string
	State string
	Since *time.Time
	To    *time.Time
}

// NewGhIssues creates a pointer to GhIssues
func NewGhIssues(owner, repo, since, to, state string) (*GhIssues, error) {
	sinceTime, err := time.Parse(time.RFC3339, since)
	if err != nil {
		return nil, fmt.Errorf("n error occured while parsing the `since` date %v: %w", since, err)
	}
	var ptrToTime *time.Time
	if to != "" {
		toTime, err := time.Parse(time.RFC3339, to)
		if err != nil {
			return nil, fmt.Errorf("n error occured while parsing the `to` date %v: %w", to, err)
		}
		ptrToTime = &toTime
	} else {
		ptrToTime = nil
	}
	ghi := GhIssues{
		Owner: owner,
		Repo:  repo,
		State: state,
		Since: &sinceTime,
		To:    ptrToTime,
	}
	return &ghi, nil

}

// GetOpts returns the github.IssueListByRepoOptions
func (ghi *GhIssues) GetOpts() *github.IssueListByRepoOptions {
	return &github.IssueListByRepoOptions{
		Sort:      "updated",
		Direction: "desc",
		State:     ghi.State,
		Since:     *ghi.Since,
	}

}

func (ghi *GhIssues) fetchPageIssue(client *github.Client, opts github.IssueListByRepoOptions, page *Page) (err error) {
	if err != nil {
		return err
	}

	opts.Page = page.Number

	ctx := context.TODO()
	issues, response, err := client.Issues.ListByRepo(ctx, ghi.Owner, ghi.Repo, &opts)
	if err == nil {
		page.Fetched = true
		page.Result = issues
	}
	page.Next = response.NextPage
	page.Last = response.LastPage
	return err
}

//IssuePager returns pages of issues
func (ghi *GhIssues) IssuePager(client *github.Client) (pager Pager, err error) {
	pager = Pager{}
	page := &Page{Number: 1}
	pager.Add(page)
	opts := ghi.GetOpts()
	err = ghi.fetchPageIssue(client, *opts, page)
	if err != nil {
		return pager, err
	}

	for i := page.Next; i <= page.Last; i++ {
		page := &Page{Number: i}
		pager.Add(page)
		go func(page *Page) {
			ghi.fetchPageIssue(client, *opts, page)
		}(page)
	}
	// Wait until all the Pages are fetched
	for !pager.IsFetched() {
		time.Sleep(1 * time.Second)
	}
	return pager, nil
}

// StringifyIssue returns a string representation of a github issue
func StringifyIssue(issue github.Issue) string {
	return fmt.Sprintf("#%d %s %s %s -- %s\n",
		issue.GetNumber(),
		issue.UpdatedAt.Format(time.RFC822),
		issue.GetState(), issue.GetUser().GetLogin(),
		issue.GetTitle(),
	)
}

func runIssues(client *github.Client, owner, repo, since, to, state string) error {
	ghi, err := NewGhIssues(owner, repo, since, to, state)
	if err != nil {
		return fmt.Errorf("runIssues: cannot create a NewGhIssues : %w", err)
	}

	pager, err := ghi.IssuePager(client)
	if err != nil {
		return fmt.Errorf("runIssues: could not query github: %w", err)
	}

	issueCount := 0

	for _, page := range pager.Pages {
		for _, issue := range page.Result {
			if ghi.To == nil || ghi.To.After(*issue.UpdatedAt) {
				issueCount++
				fmt.Printf(StringifyIssue(*issue))
			}
		}
	}
	return nil
}
