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

func (c *GitLabConnector) FindProposal(branch, target domain.LocalBranchName) (*Proposal, error) {
	opts := &gitlab.ListProjectMergeRequestsOptions{
		State:        gitlab.String("opened"),
		SourceBranch: gitlab.String(branch.String()),
		TargetBranch: gitlab.String(target.String()),
	}
	mergeRequests, _, err := c.client.MergeRequests.ListProjectMergeRequests(c.projectPath(), opts)
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

func (c *GitLabConnector) SquashMergeProposal(number int, message string) (mergeSHA domain.SHA, err error) {
	if number <= 0 {
		return domain.SHA{}, fmt.Errorf(messages.ProposalNoNumberGiven)
	}
	c.log.Start(messages.HostingGitlabMergingViaAPI, number)
	// the GitLab API wants the full commit message in the body
	result, _, err := c.client.MergeRequests.AcceptMergeRequest(c.projectPath(), number, &gitlab.AcceptMergeRequestOptions{
		SquashCommitMessage: gitlab.String(message),
		Squash:              gitlab.Bool(true),
		// the branch will be deleted by Git Town
		ShouldRemoveSourceBranch: gitlab.Bool(false),
	})
	if err != nil {
		c.log.Failed(err)
		return domain.SHA{}, err
	}
	c.log.Success()
	return domain.NewSHA(result.SHA), nil
}

func (c *GitLabConnector) UpdateProposalTarget(number int, target domain.LocalBranchName) error {
	c.log.Start(messages.HostingGitlabUpdateMRViaAPI, number, target)
	_, _, err := c.client.MergeRequests.UpdateMergeRequest(c.projectPath(), number, &gitlab.UpdateMergeRequestOptions{
		TargetBranch: gitlab.String(target.String()),
	})
	if err != nil {
		c.log.Failed(err)
		return err
	}
	c.log.Success()
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

func (c *GitLabConfig) DefaultProposalMessage(proposal Proposal) string {
	return fmt.Sprintf("%s (!%d)", proposal.Title, proposal.Number)
}

func (c *GitLabConfig) projectPath() string {
	return fmt.Sprintf("%s/%s", c.Organization, c.Repository)
}

func (c *GitLabConfig) baseURL() string {
	return fmt.Sprintf("https://%s", c.Hostname)
}

func (c *GitLabConfig) HostingServiceName() string {
	return "GitLab"
}

func (c *GitLabConfig) NewProposalURL(branch, parentBranch domain.LocalBranchName) (string, error) {
	query := url.Values{}
	query.Add("merge_request[source_branch]", branch.String())
	query.Add("merge_request[target_branch]", parentBranch.String())
	return fmt.Sprintf("%s/-/merge_requests/new?%s", c.RepositoryURL(), query.Encode()), nil
}

func (c *GitLabConfig) RepositoryURL() string {
	return fmt.Sprintf("%s/%s", c.baseURL(), c.projectPath())
}

// *************************************
// Helper functions
// *************************************

func parseGitLabMergeRequest(mergeRequest *gitlab.MergeRequest) Proposal {
	return Proposal{
		Number:          mergeRequest.IID,
		Target:          mergeRequest.TargetBranch,
		Title:           mergeRequest.Title,
		CanMergeWithAPI: true,
	}
}
