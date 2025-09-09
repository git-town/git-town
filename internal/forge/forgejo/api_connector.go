package forgejo

import (
	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

var _ forgedomain.APIConnector = forgejoConnector

func (self Connector) VerifyConnection() forgedomain.VerifyConnectionResult {
	user, _, err := self.client.GetMyUserInfo()
	if err != nil {
		return forgedomain.VerifyConnectionResult{
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
	return forgedomain.VerifyConnectionResult{
		AuthenticatedUser:   NewOption(user.UserName),
		AuthenticationError: nil,
		AuthorizationError:  err,
	}
}
