package connector

import (
	"github.com/git-town/git-town/tools/release_stats/data"
	"github.com/google/go-github/v58/github"
)

// provides the users that created the given issues
func issuesCreators(issues []*github.Issue) data.Users {
	result := data.NewUsers()
	for _, issue := range issues {
		result.AddUser(issueCreator(issue))
	}
	return result
}

func issueCreator(issue *github.Issue) string {
	return *issue.User.Login
}
