package main

import "github.com/google/go-github/v58/github"

func (gh githubConnector) issuesParticipants(issues []*github.Issue) Users {
	result := NewUsers()
	result.AddUsers(issuesCreators(issues))
	result.AddUsers(gh.issuesCommenters(issues))
	return result
}
