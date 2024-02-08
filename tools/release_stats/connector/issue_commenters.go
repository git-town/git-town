package connector

import (
	"github.com/git-town/git-town/tools/release_stats/data"
	"github.com/google/go-github/v58/github"
)

// provides the people who have commented on the given issues
func (gh Connector) issuesCommenters(issues []*github.Issue) data.Users {
	result := data.NewUsers()
	for _, issue := range issues {
		result.AddUsers(gh.issueCommenters(issue))
	}
	return result
}

func (gh Connector) issueCommenters(issue *github.Issue) data.Users {
	result := data.NewUsers()
	comments, _, err := gh.client.Issues.ListComments(gh.context, org, repo, *issue.Number, nil)
	if err != nil {
		panic(err)
	}
	for _, comment := range comments {
		result.AddUser(*comment.User.Login)
	}
	return result
}
