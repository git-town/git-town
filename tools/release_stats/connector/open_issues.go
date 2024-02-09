package connector

import (
	"fmt"

	"github.com/google/go-github/v58/github"
)

func (gh Connector) OpenIssues() []*github.Issue {
	result := []*github.Issue{}
	query := "repo:git-town/git-town is:issue is:open"
	fmt.Print("loading open issues ")
	for page := 1; ; page++ {
		results, _, err := gh.client.Search.Issues(gh.context, query, &github.SearchOptions{
			Sort:  "created",
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
