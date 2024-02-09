package connector

import (
	"fmt"

	"github.com/google/go-github/v58/github"
)

func (gh Connector) IssueReactions(issue *github.Issue) []*github.Reaction {
	result := []*github.Reaction{}
	fmt.Printf("loading reactions to #%d ", issue.GetNumber())
	for page := 0; ; page++ {
		reactions, _, err := gh.client.Reactions.ListIssueReactions(gh.context, org, repo, *issue.Number, &github.ListOptions{
			Page:    page,
			PerPage: pageSize,
		})
		if err != nil {
			panic(err.Error())
		}
		fmt.Print(".")
		if len(reactions) == 0 {
			break
		}
		result = append(result, reactions...)
	}
	fmt.Printf(" %d\n", len(result))
	return result
}
