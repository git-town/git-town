package connector

import (
	"github.com/git-town/git-town/tools/release_stats/data"
	"github.com/google/go-github/v58/github"
)

func (gh Connector) issuesCommenters(issues []*github.Issue) data.Users {
	result := data.NewUsers()
	for _, issue := range issues {
		result.AddUsers(gh.issueCommenters(issue))
	}
	return result
}

func (gh Connector) issueCommenters(issue *github.Issue) data.Users {
	return CommentsAuthors(gh.IssueComments(issue))
}
