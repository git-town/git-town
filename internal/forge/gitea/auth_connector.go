package gitea

import (
	"code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// type-check to ensure conformance to the Connector interface
var (
	giteaAuthConnector AuthConnector
	_                  forgedomain.AuthVerifier = giteaAuthConnector
	_                  forgedomain.Connector    = giteaAuthConnector
)

type AuthConnector struct {
	AnonConnector
	APIToken Option[forgedomain.GiteaToken]
	client   *gitea.Client
	log      print.Logger
}

func (self AuthConnector) VerifyConnection() forgedomain.VerifyConnectionResult {
	user, _, err := self.client.GetMyUserInfo()
	if err != nil {
		return forgedomain.VerifyConnectionResult{
			AuthenticatedUser:   None[string](),
			AuthenticationError: err,
			AuthorizationError:  nil,
		}
	}
	_, _, err = self.client.ListRepoPullRequests(self.Organization, self.Repository, gitea.ListPullRequestsOptions{
		ListOptions: gitea.ListOptions{
			PageSize: 1,
		},
	})
	return forgedomain.VerifyConnectionResult{
		AuthenticatedUser:   NewOption(user.UserName),
		AuthenticationError: nil,
		AuthorizationError:  err,
	}
}
