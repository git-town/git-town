package connector

import (
	"github.com/git-town/git-town/tools/release_stats/data"
	"github.com/google/go-github/v58/github"
)

func (gh Connector) IssuesParticipants(issues []*github.Issue) data.Users {
	result := data.NewUsers()
	for _, issue := range issues {
		result.AddUser(*issue.User.Login)
		for _, reaction := range gh.IssueReactions(issue) {
			result.AddUser(*reaction.User.Login)
		}
		for _, comment := range gh.IssueComments(issue) {
			result.AddUser(*comment.User.Login)
			for _, reaction := range gh.CommentReactions(comment) {
				result.AddUser(*reaction.User.Login)
			}
		}
	}
	return result
}
