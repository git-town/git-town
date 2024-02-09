package connector

import "github.com/google/go-github/v58/github"

func (gh Connector) IssuesComments(issues []*github.Issue) []*github.IssueComment {
	result := []*github.IssueComment{}
	for _, issue := range issues {
		result = append(result, gh.IssueComments(issue)...)
	}
	return result
}

func (gh Connector) IssueComments(issue *github.Issue) []*github.IssueComment {
	comments, _, err := gh.client.Issues.ListComments(gh.context, org, repo, *issue.Number, nil)
	if err != nil {
		panic(err)
	}
	return comments
}
