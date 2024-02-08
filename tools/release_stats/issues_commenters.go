package main

import (
	"github.com/google/go-github/v58/github"
)

// provides the people who have commented on the given issues
func (gh githubConnector) issuesCommenters(issues []*github.Issue) Users {
	result := NewUsers()
	for _, issue := range issues {
		result.AddUsers(gh.issueCommenters(issue))
	}
	return result
}

func (gh githubConnector) issueCommenters(issue *github.Issue) Users {
	result := NewUsers()
	comments, _, err := gh.client.Issues.ListComments(gh.context, org, repo, *issue.Number, nil)
	if err != nil {
		panic(err)
	}
	for _, comment := range comments {
		result.AddUser(*comment.User.Login)
	}
	return result
}
