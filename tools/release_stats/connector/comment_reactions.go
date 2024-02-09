package connector

import (
	"fmt"

	"github.com/google/go-github/v58/github"
)

func (gh Connector) CommentsReactions(comments []*github.IssueComment) []*github.Reaction {
	result := []*github.Reaction{}
	for _, comment := range comments {
		result = append(result, gh.CommentReactions(comment)...)
	}
	return result
}

func (gh Connector) CommentReactions(comment *github.IssueComment) []*github.Reaction {
	result := []*github.Reaction{}
	if *comment.Reactions.TotalCount == 0 {
		return result
	}
	fmt.Printf("loading reactions for comment %s ", *comment.URL)
	for page := 0; ; page++ {
		reactions, _, err := gh.client.Reactions.ListIssueCommentReactions(gh.context, org, repo, *comment.ID, &github.ListOptions{
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
		for _, reaction := range reactions {
			result = append(result, reaction)
		}

	}
	fmt.Printf("%d\n", len(result))
	return result
}
