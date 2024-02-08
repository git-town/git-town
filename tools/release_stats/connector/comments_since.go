package connector

import (
	"context"
	"fmt"
	"time"

	"github.com/git-town/git-town/tools/release_stats/console"
	"github.com/git-town/git-town/tools/release_stats/git"
	"github.com/google/go-github/v58/github"
)

func (gh Connector) foo() {}

// provides all users that commented anywhere since the given date
func (gh Connector) commentsSince(tag git.Tag) []*github.IssueComment {
	created := "created"
	asc := "asc"
	fmt.Printf("loading comments on issues since %s", console.Cyan.Styled(tag.ISOTime))
	result := []*github.IssueComment{}
	for page := 0; ; page++ {
		comments, _, err := gh.client.Issues.ListComments(gh.context, "git-town", "git-town", 0, &github.IssueListCommentsOptions{
			Since:     &tag.Time,
			Sort:      &created,
			Direction: &asc,
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: pageSize,
			},
		})
		if err != nil {
			panic(err)
		}
		fmt.Print(".")
		if len(comments) == 0 {
			break
		}
		for _, issueComment := range comments {
			result = append(result, issueComment)
		}
	}
	return result
}

func commentsOnPullRequestsSince(date time.Time, client *github.Client, context context.Context) []*github.PullRequestComment {
	result := []*github.PullRequestComment{}
	for page := 0; ; page++ {
		comments, _, err := client.PullRequests.ListComments(context, "git-town", "git-town", 0, &github.PullRequestListCommentsOptions{
			Since: date,
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: pageSize,
			},
		})
		if err != nil {
			panic(err)
		}
		fmt.Print(".")
		if len(comments) == 0 {
			break
		}
		for _, issueComment := range comments {
			result = append(result, issueComment)
		}
	}
	return result
}
