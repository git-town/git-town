package gitlab

import (
	"fmt"
	"net/http"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/git/giturl"
	"github.com/git-town/git-town/v11/src/hosting/hostingdomain"
	"github.com/git-town/git-town/v11/src/messages"
	"github.com/xanzy/go-gitlab"
)

// Connector provides standardized connectivity for the given repository (gitlab.com/owner/repo)
// via the GitLab API.
type Connector struct {
	client *gitlab.Client
	Config
	log hostingdomain.Log
}

func (self *Connector) FindProposal(branch, target gitdomain.LocalBranchName) (*hostingdomain.Proposal, error) {
	opts := &gitlab.ListProjectMergeRequestsOptions{
		State:        gitlab.Ptr("opened"),
		SourceBranch: gitlab.Ptr(branch.String()),
		TargetBranch: gitlab.Ptr(target.String()),
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

func (self *Connector) HostingPlatformName() string {
	return "GitLab"
}

func (self *Connector) SquashMergeProposal(number int, message string) (mergeSHA gitdomain.SHA, err error) {
	if number <= 0 {
		return gitdomain.EmptySHA(), fmt.Errorf(messages.ProposalNoNumberGiven)
	}
	self.log.Start(messages.HostingGitlabMergingViaAPI, number)
	// the GitLab API wants the full commit message in the body
	result, _, err := self.client.MergeRequests.AcceptMergeRequest(self.projectPath(), number, &gitlab.AcceptMergeRequestOptions{
		SquashCommitMessage: gitlab.Ptr(message),
		Squash:              gitlab.Ptr(true),
		// the branch will be deleted by Git Town
		ShouldRemoveSourceBranch: gitlab.Ptr(false),
	})
	if err != nil {
		self.log.Failed(err)
		return gitdomain.EmptySHA(), err
	}
	self.log.Success()
	return gitdomain.NewSHA(result.SHA), nil
}

func (self *Connector) UpdateProposalTarget(number int, target gitdomain.LocalBranchName) error {
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
func NewConnector(args NewConnectorArgs) (*Connector, error) {
	if args.OriginURL == nil || (args.OriginURL.Host != "gitlab.com" && args.HostingPlatform != configdomain.HostingPlatformGitLab) {
		return nil, nil //nolint:nilnil
	}
	gitlabConfig := Config{
		Config: hostingdomain.Config{
			Hostname:     args.OriginURL.Host,
			Organization: args.OriginURL.Org,
			Repository:   args.OriginURL.Repo,
		},
		APIToken: args.APIToken,
	}
	clientOptFunc := gitlab.WithBaseURL(gitlabConfig.baseURL())
	httpClient := gitlab.WithHTTPClient(&http.Client{})
	client, err := gitlab.NewOAuthClient(gitlabConfig.APIToken.String(), httpClient, clientOptFunc)
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
	HostingPlatform configdomain.HostingPlatform
	OriginURL       *giturl.Parts
	APIToken        configdomain.GitLabToken
	Log             hostingdomain.Log
}

func parseMergeRequest(mergeRequest *gitlab.MergeRequest) hostingdomain.Proposal {
	return hostingdomain.Proposal{
		Number:       mergeRequest.IID,
		Target:       gitdomain.NewLocalBranchName(mergeRequest.TargetBranch),
		Title:        mergeRequest.Title,
		MergeWithAPI: true,
	}
}
