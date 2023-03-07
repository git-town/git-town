package hosting

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/git-town/git-town/v7/src/config"
	"github.com/xanzy/go-gitlab"
)

// GitLabConnector provides standardized connectivity for the given repository (gitlab.com/owner/repo)
// via the GitLab API.
type GitLabConnector struct {
	client *gitlab.Client
	GitLabConfig
	log logFn
}

func (c *GitLabConnector) FindProposal(branch, target string) (*Proposal, error) {
	opts := &gitlab.ListProjectMergeRequestsOptions{
		State:        gitlab.String("opened"),
		SourceBranch: gitlab.String(branch),
		TargetBranch: gitlab.String(target),
	}
	mergeRequests, _, err := c.client.MergeRequests.ListProjectMergeRequests(c.projectPath(), opts)
	if err != nil {
		return nil, err
	}
	if len(mergeRequests) == 0 {
		return nil, nil //nolint:nilnil
	}
	if len(mergeRequests) > 1 {
		return nil, fmt.Errorf("found %d merge requests for branch %q", len(mergeRequests), branch)
	}
	proposal := parseGitLabMergeRequest(mergeRequests[0])
	return &proposal, nil
}

//nolint:nonamedreturns  // return value isn't obvious from function name
func (c *GitLabConnector) SquashMergeProposal(number int, message string) (mergeSHA string, err error) {
	if number <= 0 {
		return "", fmt.Errorf("no merge request number given")
	}
	if c.log != nil {
		c.log("GitLab API: Merging MR !%d\n", number)
	}
	// the GitLab API wants the full commit message in the body
	result, _, err := c.client.MergeRequests.AcceptMergeRequest(c.projectPath(), number, &gitlab.AcceptMergeRequestOptions{
		SquashCommitMessage: gitlab.String(message),
		Squash:              gitlab.Bool(true),
		// the branch will be deleted by Git Town
		ShouldRemoveSourceBranch: gitlab.Bool(false),
	})
	if err != nil {
		return "", err
	}
	return result.SHA, nil
}

func (c *GitLabConnector) UpdateProposalTarget(number int, target string) error {
	if c.log != nil {
		c.log("GitLab API: Updating target branch for MR !%d to %q\n", number, target)
	}
	_, _, err := c.client.MergeRequests.UpdateMergeRequest(c.projectPath(), number, &gitlab.UpdateMergeRequestOptions{
		TargetBranch: gitlab.String(target),
	})
	return err
}

// NewGitlabConfig provides GitLab configuration data if the current repo is hosted on GitLab,
// otherwise nil.
func NewGitlabConnector(gitConfig gitTownConfig, log logFn) (*GitLabConnector, error) {
	hostingService, err := gitConfig.HostingService()
	if err != nil {
		return nil, err
	}
	url := gitConfig.OriginURL()
	if url == nil || (url.Host != "gitlab.com" && hostingService != config.HostingServiceGitLab) {
		return nil, nil //nolint:nilnil
	}
	gitlabConfig := GitLabConfig{CommonConfig{
		APIToken:     gitConfig.GitLabToken(),
		Hostname:     url.Host,
		Organization: url.Org,
		Repository:   url.Repo,
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
		log:          log,
	}
	return &connector, nil
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

func (c *GitLabConfig) NewProposalURL(branch, parentBranch string) (string, error) {
	query := url.Values{}
	query.Add("merge_request[source_branch]", branch)
	query.Add("merge_request[target_branch]", parentBranch)
	return fmt.Sprintf("%s/merge_requests/new?%s", c.RepositoryURL(), query.Encode()), nil
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
