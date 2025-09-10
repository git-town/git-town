package bitbucketcloud

import (
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/ktrysmt/go-bitbucket"
)

var (
	bbclAuthConnector AuthConnector
	_                 forgedomain.AuthVerifier = bbclAuthConnector
	_                 forgedomain.Connector    = bbclAuthConnector
)

// AuthConnector provides access to the Bitbucket Cloud API.
type AuthConnector struct {
	AnonConnector
	client Mutable[bitbucket.Client]
	log    print.Logger
}

func (self AuthConnector) VerifyCredentials() forgedomain.VerifyCredentialsResult {
	user, err := self.client.Value.User.Profile()
	if err != nil {
		return forgedomain.VerifyCredentialsResult{
			AuthenticatedUser:   None[string](),
			AuthenticationError: err,
			AuthorizationError:  nil,
		}
	}
	_, err = self.client.Value.Repositories.PullRequests.Gets(&bitbucket.PullRequestsOptions{
		Owner:    self.Organization,
		RepoSlug: self.Repository,
		Query:    "",
		States:   []string{},
	})
	return forgedomain.VerifyCredentialsResult{
		AuthenticatedUser:   NewOption(user.Username),
		AuthenticationError: nil,
		AuthorizationError:  err,
	}
}
