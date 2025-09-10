package gitlab

import (
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
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
	client   *gitlab.Client
	APIToken forgedomain.GitLabToken
	log      print.Logger
}

func (self AuthConnector) VerifyConnection() forgedomain.VerifyConnectionResult {
	user, _, err := self.client.Users.CurrentUser()
	if err != nil {
		return forgedomain.VerifyConnectionResult{
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
	return forgedomain.VerifyConnectionResult{
		AuthenticatedUser:   NewOption(user.Username),
		AuthenticationError: nil,
		AuthorizationError:  err,
	}
}

func parseMergeRequest(mergeRequest *gitlab.BasicMergeRequest) forgedomain.ProposalData {
	return forgedomain.ProposalData{
		MergeWithAPI: true,
		Number:       mergeRequest.IID,
		Source:       gitdomain.NewLocalBranchName(mergeRequest.SourceBranch),
		Target:       gitdomain.NewLocalBranchName(mergeRequest.TargetBranch),
		Title:        mergeRequest.Title,
		Body:         NewOption(mergeRequest.Description),
		URL:          mergeRequest.WebURL,
	}
}
