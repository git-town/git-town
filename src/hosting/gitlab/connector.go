package gitlab

import (
	"fmt"
	"net/http"

	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/domain"
	"github.com/git-town/git-town/v10/src/git/giturl"
	"github.com/git-town/git-town/v10/src/hosting/common"
	"github.com/git-town/git-town/v10/src/messages"
	"github.com/xanzy/go-gitlab"
)

// Connector provides standardized connectivity for the given repository (gitlab.com/owner/repo)
// via the GitLab API.
type Connector struct {
	client *gitlab.Client
	Config
	log common.Log
}

func (self *Connector) FindProposal(branch, target domain.LocalBranchName) (*domain.Proposal, error) {
	opts := &gitlab.ListProjectMergeRequestsOptions{
		State:        gitlab.String("opened"),
		SourceBranch: gitlab.String(branch.String()),
		TargetBranch: gitlab.String(target.String()),
	}
	mergeRequests, _, err := self.client.MergeRequests.ListProjectMergeRequests(self.projectPath(), opts)
	if err != nil {
		return nil, err
	}
	if len(mergeRequests) == 0 {
		return nil, nil //nolint:nilnil
	}
	if len(mergeRequests) > 1 {
		return nil, fmt.Errorf(messages.ProposalMultipleFound, len(mergeRequests), branch, target)
	}
	proposal := parseMergeRequest(mergeRequests[0])
	return &proposal, nil
}

func (self *Connector) SquashMergeProposal(number int, message string) (mergeSHA domain.SHA, err error) {
	if number <= 0 {
		return domain.EmptySHA(), fmt.Errorf(messages.ProposalNoNumberGiven)
	}
	self.log.Start(messages.HostingGitlabMergingViaAPI, number)
	// the GitLab API wants the full commit message in the body
	result, _, err := self.client.MergeRequests.AcceptMergeRequest(self.projectPath(), number, &gitlab.AcceptMergeRequestOptions{
		SquashCommitMessage: gitlab.String(message),
		Squash:              gitlab.Bool(true),
		// the branch will be deleted by Git Town
		ShouldRemoveSourceBranch: gitlab.Bool(false),
	})
	if err != nil {
		self.log.Failed(err)
		return domain.EmptySHA(), err
	}
	self.log.Success()
	return domain.NewSHA(result.SHA), nil
}

func (self *Connector) UpdateProposalTarget(number int, target domain.LocalBranchName) error {
	self.log.Start(messages.HostingGitlabUpdateMRViaAPI, number, target)
	_, _, err := self.client.MergeRequests.UpdateMergeRequest(self.projectPath(), number, &gitlab.UpdateMergeRequestOptions{
		TargetBranch: gitlab.String(target.String()),
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
func NewConnector(args NewConnectorArgs) (*Connector, error) {
	if args.OriginURL == nil || (args.OriginURL.Host != "gitlab.com" && args.HostingService != config.HostingGitLab) {
		return nil, nil //nolint:nilnil
	}
	gitlabConfig := Config{common.Config{
		APIToken:     args.APIToken,
		Hostname:     args.OriginURL.Host,
		Organization: args.OriginURL.Org,
		Repository:   args.OriginURL.Repo,
	}}
	clientOptFunc := gitlab.WithBaseURL(gitlabConfig.baseURL())
	httpClient := gitlab.WithHTTPClient(&http.Client{})
	client, err := gitlab.NewOAuthClient(gitlabConfig.APIToken, httpClient, clientOptFunc)
	if err != nil {
		return nil, err
	}
	connector := Connector{
		client: client,
		Config: gitlabConfig,
		log:    args.Log,
	}
	return &connector, nil
}

type NewConnectorArgs struct {
	HostingService config.Hosting
	OriginURL      *giturl.Parts
	APIToken       string
	Log            common.Log
}

func parseMergeRequest(mergeRequest *gitlab.MergeRequest) domain.Proposal {
	return domain.Proposal{
		Number:       mergeRequest.IID,
		Target:       domain.NewLocalBranchName(mergeRequest.TargetBranch),
		Title:        mergeRequest.Title,
		MergeWithAPI: true,
	}
}
