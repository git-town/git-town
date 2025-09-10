package github

import (
	"context"

	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/google/go-github/v58/github"
)

// type-check to ensure conformance to the Connector interface
var (
	githubAuthConnector AuthConnector
	_                   forgedomain.AuthVerifier = githubAuthConnector
	_                   forgedomain.Connector    = githubAuthConnector
)

// Connector provides standardized connectivity for the given repository (github.com/owner/repo)
// via the GitHub API.
type AuthConnector struct {
	AnonConnector
	APIToken Option[forgedomain.GitHubToken]
	client   Mutable[github.Client]
	log      print.Logger
}

func (self AuthConnector) VerifyConnection() forgedomain.VerifyConnectionResult {
	user, _, err := self.client.Value.Users.Get(context.Background(), "")
	if err != nil {
		return forgedomain.VerifyConnectionResult{
			AuthenticatedUser:   None[string](),
			AuthenticationError: err,
			AuthorizationError:  nil,
		}
	}
	_, _, err = self.client.Value.PullRequests.List(context.Background(), self.Organization, self.Repository, &github.PullRequestListOptions{
		ListOptions: github.ListOptions{
			PerPage: 1,
		},
	})
	return forgedomain.VerifyConnectionResult{
		AuthenticatedUser:   NewOption(*user.Login),
		AuthenticationError: nil,
		AuthorizationError:  err,
	}
}
