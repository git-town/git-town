package gitlab

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/git-town/git-town/v15/internal/cli/print"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/git/giturl"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
	"github.com/git-town/git-town/v15/internal/hosting/hostingdomain"
	"github.com/git-town/git-town/v15/internal/messages"
	"github.com/xanzy/go-gitlab"
)

// Connector provides standardized connectivity for the given repository (gitlab.com/owner/repo)
// via the GitLab API.
type Connector struct {
	client *gitlab.Client
	Data
	log print.Logger
}

func (self Connector) FindProposal(branch, target gitdomain.LocalBranchName) (Option[hostingdomain.Proposal], error) {
	opts := &gitlab.ListProjectMergeRequestsOptions{
		State:        gitlab.Ptr("opened"),
		SourceBranch: gitlab.Ptr(branch.String()),
		TargetBranch: gitlab.Ptr(target.String()),
	}
	mergeRequests, _, err := self.client.MergeRequests.ListProjectMergeRequests(self.projectPath(), opts)
	if err != nil {
		return None[hostingdomain.Proposal](), err
	}
	switch len(mergeRequests) {
	case 0:
		return None[hostingdomain.Proposal](), nil
	case 1:
		return Some(parseMergeRequest(mergeRequests[0])), nil
	default:
		return None[hostingdomain.Proposal](), fmt.Errorf(messages.ProposalMultipleFound, len(mergeRequests), branch, target)
	}
}

func (self Connector) SquashMergeProposal(number int, message gitdomain.CommitMessage) error {
	if number <= 0 {
		return errors.New(messages.ProposalNoNumberGiven)
	}
	self.log.Start(messages.HostingGitlabMergingViaAPI, number)
	// the GitLab API wants the full commit message in the body
	_, _, err := self.client.MergeRequests.AcceptMergeRequest(self.projectPath(), number, &gitlab.AcceptMergeRequestOptions{
		SquashCommitMessage: gitlab.Ptr(message.String()),
		Squash:              gitlab.Ptr(true),
		// the branch will be deleted by Git Town
		ShouldRemoveSourceBranch: gitlab.Ptr(false),
	})
	if err != nil {
		self.log.Failed(err)
		return err
	}
	self.log.Success()
	return nil
}

func (self Connector) UpdateProposalTarget(number int, target gitdomain.LocalBranchName) error {
	self.log.Start(messages.HostingGitlabUpdateMRViaAPI, number, target)
	_, _, err := self.client.MergeRequests.UpdateMergeRequest(self.projectPath(), number, &gitlab.UpdateMergeRequestOptions{
		TargetBranch: gitlab.Ptr(target.String()),
	})
	if err != nil {
		self.log.Failed(err)
		return err
	}
	self.log.Success()
	return nil
}

// NewGitlabConfig provides GitLab configuration data if the current repo is hosted on GitLab,
// otherwise nil.
func NewConnector(args NewConnectorArgs) (Connector, error) {
	gitlabData := Data{
		APIToken: args.APIToken,
		Data: hostingdomain.Data{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
	}
	clientOptFunc := gitlab.WithBaseURL(gitlabData.baseURL())
	httpClient := gitlab.WithHTTPClient(&http.Client{}) //exhaustruct:ignore
	client, err := gitlab.NewOAuthClient(gitlabData.APIToken.String(), httpClient, clientOptFunc)
	if err != nil {
		return Connector{}, err
	}
	connector := Connector{
		Data:   gitlabData,
		client: client,
		log:    args.Log,
	}
	return connector, nil
}

type NewConnectorArgs struct {
	APIToken  Option[configdomain.GitLabToken]
	Log       print.Logger
	RemoteURL giturl.Parts
}

func parseMergeRequest(mergeRequest *gitlab.MergeRequest) hostingdomain.Proposal {
	return hostingdomain.Proposal{
		Number:       mergeRequest.IID,
		Target:       gitdomain.NewLocalBranchName(mergeRequest.TargetBranch),
		Title:        mergeRequest.Title,
		MergeWithAPI: true,
	}
}
