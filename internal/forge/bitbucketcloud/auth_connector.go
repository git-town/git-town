package bitbucketcloud

import (
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/ktrysmt/go-bitbucket"
)

var (
	bbclAPIConnector AuthConnector
	_                forgedomain.AuthVerifier = bbclAPIConnector
	_                forgedomain.Connector    = bbclAPIConnector
)

// AuthConnector implements the connector functionality if API credentials are available.
type AuthConnector struct {
	AnonConnector
	client *bitbucket.Client
	log    print.Logger
}

func (self AuthConnector) VerifyCredentials() forgedomain.VerifyCredentialsResult {
	user, err := self.client.User.Profile()
	if err != nil {
		return forgedomain.VerifyCredentialsResult{
			AuthenticatedUser:   None[string](),
			AuthenticationError: err,
			AuthorizationError:  nil,
		}
	}
	_, err = self.client.Repositories.PullRequests.Gets(&bitbucket.PullRequestsOptions{
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
