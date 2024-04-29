package connector

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/tools/stats_release/console"
	"github.com/google/go-github/v58/github"
)

func (gh Connector) IssueReactions(issue *github.Issue, issueType string, current, total int) []*github.Reaction {
	result := []*github.Reaction{}
	fmt.Printf("loading reactions to %s %d/%d (#%d) ", issueType, current, total, issue.GetNumber())
	for page := 1; ; page++ {
		reactions, response, err := gh.client.Reactions.ListIssueReactions(gh.context, org, repo, *issue.Number, &github.ListOptions{
			Page:    page,
			PerPage: pageSize,
		})
		if err != nil {
			panic(err.Error())
		}
		result = append(result, reactions...)
		if response.NextPage == 0 {
			break
		}
	}
	fmt.Println(console.Green.Styled(" ok"))
	if len(result) > 0 {
		texts := make([]string, len(result))
		for r, reaction := range result {
			texts[r] = fmt.Sprintf("%s (%s)", *reaction.Content, reactionAuthor(*reaction))
		}
		fmt.Printf("  %s\n", strings.Join(texts, ", "))
	}
	return result
}
