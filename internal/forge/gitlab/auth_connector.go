package gitlab

import (
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// type-check to ensure conformance to the Connector interface
var (
	gitlabAuthConnector AuthConnector
	_                   forgedomain.AuthVerifier = gitlabAuthConnector
	_                   forgedomain.Connector    = gitlabAuthConnector
)

// Connector provides standardized connectivity for the given repository (gitlab.com/owner/repo)
// via the GitLab API.
type AuthConnector struct {
	AnonConnector
	APIToken forgedomain.GitLabToken
	client   *gitlab.Client
	log      print.Logger
}

func (self AuthConnector) VerifyCredentials() forgedomain.VerifyCredentialsResult {
	user, _, err := self.client.Users.CurrentUser()
	if err != nil {
		return forgedomain.VerifyCredentialsResult{
			AuthenticatedUser:   None[string](),
			AuthenticationError: err,
			AuthorizationError:  nil,
		}
	}
	_, _, err = self.client.MergeRequests.ListMergeRequests(&gitlab.ListMergeRequestsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 1,
		},
	})
	return forgedomain.VerifyCredentialsResult{
		AuthenticatedUser:   NewOption(user.Username),
		AuthenticationError: nil,
		AuthorizationError:  err,
	}
}
