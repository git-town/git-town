package hosting

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/git-town/git-town/v7/src/giturl"
	"github.com/xanzy/go-gitlab"
)

// GitLabConnector provides standardized connectivity for the given repository (gitlab.com/owner/repo)
// via the GitLab API.
type GitLabConnector struct {
	client *gitlab.Client
	GitLabConfig
	log logFn
}

func (c *GitLabConnector) ProposalForBranch(branch string) (*Proposal, error) {
	opts := &gitlab.ListProjectMergeRequestsOptions{
		State:        gitlab.String("opened"),
		SourceBranch: gitlab.String(branch),
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
	changeRequest := parseMergeRequest(mergeRequests[0])
	return &changeRequest, nil
}

//nolint:nonamedreturns  // return value isn't obvious from function name
func (c *GitLabConnector) SquashMergeChangeRequest(number int, message string) (mergeSHA string, err error) {
	// TODO: update PR target? Probably better to check the target here,
	// warn if it is different on GitLab than it is locally,
	// and update and merge only if a "--force" option is given.
	if number <= 0 {
		return "", fmt.Errorf("cannot merge via GitLab since there is no merge request")
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

func (c *GitLabConnector) UpdateChangeRequestTarget(number int, target string) error {
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
func NewGitlabConnector(url giturl.Parts, config gitConfig, log logFn) (*GitLabConnector, error) {
	manualHostName := config.OriginOverride()
	if manualHostName != "" {
		url.Host = manualHostName
	}
	if config.HostingService() != "gitlab" && url.Host != "gitlab.com" {
		return nil, nil //nolint:nilnil
	}
	gitlabConfig := GitLabConfig{Config{
		apiToken:   config.GitLabToken(),
		originURL:  config.OriginURL(),
		hostname:   url.Host,
		owner:      url.Org,
		repository: url.Repo,
	}}
	clientOptFunc := gitlab.WithBaseURL(gitlabConfig.baseURL())
	httpClient := gitlab.WithHTTPClient(&http.Client{})
	client, err := gitlab.NewOAuthClient(gitlabConfig.apiToken, httpClient, clientOptFunc)
	if err != nil {
		return nil, err
	}
	driver := GitLabConnector{
		client:       client,
		GitLabConfig: gitlabConfig,
		log:          log,
	}
	return &driver, nil
}

// *************************************
// GitLabConfig
// *************************************

type GitLabConfig struct {
	Config
}

func (c *GitLabConfig) DefaultCommitMessage(changeRequest Proposal) string {
	return fmt.Sprintf("%s (!%d)", changeRequest.Title, changeRequest.Number)
}

func (c *GitLabConfig) projectPath() string {
	return fmt.Sprintf("%s/%s", c.owner, c.repository)
}

func (c *GitLabConfig) baseURL() string {
	return fmt.Sprintf("https://%s", c.hostname)
}

func (c *GitLabConfig) HostingServiceName() string {
	return "GitLab"
}

func (c *GitLabConfig) NewChangeRequestURL(branch, parentBranch string) (string, error) {
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

func parseMergeRequest(mergeRequest *gitlab.MergeRequest) Proposal {
	return Proposal{
		Number:          mergeRequest.IID,
		Title:           mergeRequest.Title,
		CanMergeWithAPI: true,
	}
}
