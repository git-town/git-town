package connector

import (
	"github.com/git-town/git-town/v22/pkg/asserts"
	"github.com/google/go-github/v58/github"
)

func (gh Connector) IssueComments(issue *github.Issue) []*github.IssueComment {
	comments, _ := asserts.NoError2(gh.client.Issues.ListComments(gh.context, org, repo, *issue.Number, nil))
	return comments
}
