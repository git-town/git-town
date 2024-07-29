package connector

import (
	"github.com/git-town/git-town/tools/stats_release/data"
	"github.com/google/go-github/v58/github"
)

func (gh Connector) IssuesParticipants(issues []*github.Issue, issueType string) data.Users {
	result := data.Users{}
	total := len(issues)
	for i, issue := range issues {
		result.Add(issueAuthor(*issue))
		for _, reaction := range gh.IssueReactions(issue, issueType, i+1, total) {
			result.Add(reactionAuthor(*reaction))
		}
		for _, comment := range gh.IssueComments(issue) {
			result.Add(commentAuthor(*comment))
			for _, reaction := range gh.CommentReactions(*comment) {
				result.Add(reactionAuthor(*reaction))
			}
		}
	}
	return result
}
