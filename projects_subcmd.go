package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/github"
)

type ProjectCards struct {
	Owner    string
	Repo     string
	ColumnID int64
}

// NewProjectCards returns a pointer to the ProjectCards
func NewProjectCards(owner, repo string, columnID int64) *ProjectCards {
	return &ProjectCards{owner, repo, columnID}
}

// GetOpts return the options
func (pc *ProjectCards) GetOpts() *github.ProjectCardListOptions {
	return &github.ProjectCardListOptions{
		ArchivedState: github.String("all"),
		ListOptions:   github.ListOptions{Page: 1},
	}
}

// FetchCards returns the Cards and and error.
func (pc *ProjectCards) FetchCards(client *github.Client) ([]*github.ProjectCard, error) {
	opts := pc.GetOpts()
	ctx := context.TODO()
	var allCards []*github.ProjectCard
	for {
		cards, resp, err := client.Projects.ListProjectCards(ctx, pc.ColumnID, opts)
		if err != nil {
			return nil, fmt.Errorf("can not fetch cards: %w", err)
		}
		allCards = append(allCards, cards...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return allCards, nil

}

// StringifyCard retuns the string representation of the cards
func StringifyCard(card github.ProjectCard) string {
	if card.ContentURL != nil {
		return fmt.Sprintf("%s, last updated: %s", *card.ContentURL, card.UpdatedAt.Format(time.RFC822))
	}
	return fmt.Sprintf("Note: %s, last updated: %s", *card.Note, card.UpdatedAt.Format(time.RFC822))
}

func reportCards(client *github.Client, owner, repo string, columnID int64) error {
	pc := NewProjectCards(owner, repo, columnID)
	cards, err := pc.FetchCards(client)
	if err != nil {
		return fmt.Errorf("reportCards %w", err)
	}
	for _, card := range cards {
		fmt.Println(StringifyCard(*card))

	}
	return nil
}
