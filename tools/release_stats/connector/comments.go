package connector

import "github.com/google/go-github/v58/github"

func (gh Connector) AllComments() []*github.IssueComment {
	comments, _, err := gh.client.Issues.ListComments(gh.context, org, repo, 0, nil)
	if err != nil {
		panic(err)
	}
	return comments
}
