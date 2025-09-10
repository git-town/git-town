package forgejo

import (
	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// type-check to ensure conformance to the Connector interface
var (
	authConnector AuthConnector
	_             forgedomain.AuthVerifier = authConnector
	_             forgedomain.Connector    = authConnector
)

// AuthConnector provides access to the Forgejo API.
type AuthConnector struct {
	WebConnector
	APIToken Option[forgedomain.ForgejoToken]
	client   *forgejo.Client
	log      print.Logger
}

func (self AuthConnector) VerifyCredentials() forgedomain.VerifyCredentialsResult {
	user, _, err := self.client.GetMyUserInfo()
	if err != nil {
		return forgedomain.VerifyCredentialsResult{
			AuthenticatedUser:   None[string](),
			AuthenticationError: err,
			AuthorizationError:  nil,
		}
	}
	_, _, err = self.client.ListRepoPullRequests(self.Organization, self.Repository, forgejo.ListPullRequestsOptions{
		ListOptions: forgejo.ListOptions{
			PageSize: 1,
		},
	})
	return forgedomain.VerifyCredentialsResult{
		AuthenticatedUser:   NewOption(user.UserName),
		AuthenticationError: nil,
		AuthorizationError:  err,
	}
}
