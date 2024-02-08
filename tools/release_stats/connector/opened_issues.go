package connector

import (
	"fmt"

	"github.com/google/go-github/v58/github"
)

func (gh Connector) OpenedIssuesOrPRsSince(date string) []*github.Issue {
	// load all issues opened since the last release
	result := []*github.Issue{}
	query := fmt.Sprintf("repo:git-town/git-town opened:>=%s", date)
	fmt.Print("loading opened issues ")
	for page := 0; ; page++ {
		results, _, err := gh.client.Search.Issues(gh.context, query, &github.SearchOptions{
			Sort:  "opened",
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
			result = append(result, issue)
		}
	}
	return result
}
