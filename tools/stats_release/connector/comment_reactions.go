package connector

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/tools/stats_release/console"
	"github.com/google/go-github/v58/github"
)

func (gh Connector) CommentReactions(comment github.IssueComment) []*github.Reaction {
	result := []*github.Reaction{}
	if *comment.Reactions.TotalCount == 0 {
		return result
	}
	fmt.Printf("loading reactions to comment #%d ", comment.GetID())
	for page := 1; ; page++ {
		reactions, response, err := gh.client.Reactions.ListIssueCommentReactions(gh.context, org, repo, *comment.ID, &github.ListOptions{
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
	if len(result) == 0 {
		fmt.Println(console.Green.Styled(" ok"))
	} else {
		texts := make([]string, len(result))
		for r, reaction := range result {
			texts[r] = fmt.Sprintf("%s (%s)", *reaction.Content, reactionAuthor(*reaction))
		}
		fmt.Printf(" %s\n  %s\n", console.Green.Styled("ok"), strings.Join(texts, ", "))
	}
	return result
}
