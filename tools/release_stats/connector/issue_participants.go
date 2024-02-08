package connector

import (
	"github.com/git-town/git-town/tools/release_stats/data"
	"github.com/google/go-github/v58/github"
)

func (gh Connector) IssuesParticipants(issues []*github.Issue) data.Users {
	result := data.NewUsers()
	result.AddUsers(issuesCreators(issues))
	result.AddUsers(gh.issuesCommenters(issues))
	return result
}
