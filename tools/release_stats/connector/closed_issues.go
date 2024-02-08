package connector

import (
	"fmt"

	"github.com/google/go-github/v58/github"
)

// loads all issues and pull requests closed since the given date
func (gh Connector) ClosedIssues(date string) (closedIssues []*github.Issue, closedPullRequests []*github.Issue) {
	query := fmt.Sprintf("repo:git-town/git-town closed:>=%s", date)
	fmt.Print("loading closed issues and pull requests ")
	for page := 0; ; page++ {
		results, _, err := gh.client.Search.Issues(gh.context, query, &github.SearchOptions{
			Sort:  "closed",
			Order: "asc",
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: pageSize,
			},
		})
		if err != nil {
			panic(err)
		}
		fmt.Print(".")
		if len(results.Issues) == 0 {
			break
		}
		for _, issue := range results.Issues {
			if issue.IsPullRequest() {
				closedPullRequests = append(closedPullRequests, issue)
			} else {
				closedIssues = append(closedIssues, issue)
			}
		}
	}
	return
}
