package bitbucketcloud

import (
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/ktrysmt/go-bitbucket"
)

var bbclAPIConnector APIConnector
var _ forgedomain.APIConnector = bbclAPIConnector
var _ forgedomain.Connector = bbclAPIConnector

// APIConnector implements the connector functionality if API credentials are available.
type APIConnector struct {
	WebConnector
	client *bitbucket.Client
	log    print.Logger
}

func (self APIConnector) VerifyConnection() forgedomain.VerifyConnectionResult {
	user, err := self.client.User.Profile()
	if err != nil {
		return forgedomain.VerifyConnectionResult{
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
	return forgedomain.VerifyConnectionResult{
		AuthenticatedUser:   NewOption(user.Username),
		AuthenticationError: nil,
		AuthorizationError:  err,
	}
}
