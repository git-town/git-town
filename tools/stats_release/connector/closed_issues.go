package connector

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/tools/stats_release/console"
	"github.com/git-town/git-town/v22/pkg/asserts"
	"github.com/google/go-github/v58/github"
)

// loads all issues and pull requests closed since the given date
func (gh Connector) ClosedIssues(date string) ClosedResult {
	query := fmt.Sprintf("repo:%s/%s closed:>=%s", org, repo, date)
	fmt.Printf("loading issues and pull requests closed since %s ", date)
	closedIssues := []*github.Issue{}
	closedPullRequests := []*github.Issue{}
	for page := 1; ; page++ {
		results, response := asserts.NoError2(gh.client.Search.Issues(gh.context, query, &github.SearchOptions{
			Sort:  "closed",
			Order: "asc",
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: pageSize,
			},
		}))
		fmt.Print(".")
		for _, issue := range results.Issues {
			if issue.IsPullRequest() {
				closedPullRequests = append(closedPullRequests, issue)
			} else {
				closedIssues = append(closedIssues, issue)
			}
		}
		if response.NextPage == 0 {
			break
		}
	}
	fmt.Printf(" %s issues, %s pull requests\n", console.Green.Styled(strconv.Itoa(len(closedIssues))), console.Green.Styled(strconv.Itoa(len(closedPullRequests))))
	return ClosedResult{
		Issues:       closedIssues,
		PullRequests: closedPullRequests,
	}
}

// ClosedResult describes the issues and pull requests that are closed since the given date.
type ClosedResult struct {
	Issues       []*github.Issue
	PullRequests []*github.Issue
}
