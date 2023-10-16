package hosting

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/giturl"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/xanzy/go-gitlab"
)

// GitLabConnector provides standardized connectivity for the given repository (gitlab.com/owner/repo)
// via the GitLab API.
type GitLabConnector struct {
	client *gitlab.Client
	GitLabConfig
	log Log
}

func (self *GitLabConnector) FindProposal(branch, target domain.LocalBranchName) (*Proposal, error) {
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
	proposal := parseGitLabMergeRequest(mergeRequests[0])
	return &proposal, nil
}

func (self *GitLabConnector) SquashMergeProposal(number int, message string) (mergeSHA domain.SHA, err error) {
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

func (self *GitLabConnector) UpdateProposalTarget(number int, target domain.LocalBranchName) error {
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
func NewGitlabConnector(args NewGitlabConnectorArgs) (*GitLabConnector, error) {
	if args.OriginURL == nil || (args.OriginURL.Host != "gitlab.com" && args.HostingService != config.HostingGitLab) {
		return nil, nil //nolint:nilnil
	}
	gitlabConfig := GitLabConfig{CommonConfig{
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
	connector := GitLabConnector{
		client:       client,
		GitLabConfig: gitlabConfig,
		log:          args.Log,
	}
	return &connector, nil
}

type NewGitlabConnectorArgs struct {
	HostingService config.Hosting
	OriginURL      *giturl.Parts
	APIToken       string
	Log            Log
}

// *************************************
// GitLabConfig
// *************************************

type GitLabConfig struct {
	CommonConfig
}

func (self *GitLabConfig) DefaultProposalMessage(proposal Proposal) string {
	return fmt.Sprintf("%s (!%d)", proposal.Title, proposal.Number)
}

func (self *GitLabConfig) projectPath() string {
	return fmt.Sprintf("%s/%s", self.Organization, self.Repository)
}

func (self *GitLabConfig) baseURL() string {
	return fmt.Sprintf("https://%s", self.Hostname)
}

func (self *GitLabConfig) HostingServiceName() string {
	return "GitLab"
}

func (self *GitLabConfig) NewProposalURL(branch, parentBranch domain.LocalBranchName) (string, error) {
	query := url.Values{}
	query.Add("merge_request[source_branch]", branch.String())
	query.Add("merge_request[target_branch]", parentBranch.String())
	return fmt.Sprintf("%s/-/merge_requests/new?%s", self.RepositoryURL(), query.Encode()), nil
}

func (self *GitLabConfig) RepositoryURL() string {
	return fmt.Sprintf("%s/%s", self.baseURL(), self.projectPath())
}

// *************************************
// Helper functions
// *************************************

func parseGitLabMergeRequest(mergeRequest *gitlab.MergeRequest) Proposal {
	return Proposal{
		Number:          mergeRequest.IID,
		Target:          domain.NewLocalBranchName(mergeRequest.TargetBranch),
		Title:           mergeRequest.Title,
		CanMergeWithAPI: true,
	}
}
