package hosting

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/git-town/git-town/v7/src/giturl"
	"github.com/xanzy/go-gitlab"
)

// GitlabConfig contains connection information for GitLab based hosting platforms.
type GitlabConfig struct {
	apiToken   string
	hostname   string
	originURL  string
	owner      string
	repository string
}

// NewGitlabConfig provides GitLab configuration data if the current repo is hosted on GitLab,
// otherwise nil.
func NewGitlabConfig(url giturl.Parts, config config) *GitlabConfig {
	driverType := config.HostingService()
	manualHostName := config.OriginOverride()
	if manualHostName != "" {
		url.Host = manualHostName
	}
	if driverType != "gitlab" && url.Host != "gitlab.com" {
		return nil
	}
	return &GitlabConfig{
		apiToken:   config.GitLabToken(),
		originURL:  config.OriginURL(),
		hostname:   url.Host,
		owner:      url.Org,
		repository: url.Repo,
	}
}

func (c GitlabConfig) BaseURL() string {
	return fmt.Sprintf("https://%s", c.hostname)
}

func (c GitlabConfig) defaultCommitMessage(mergeRequest *gitlab.MergeRequest) string {
	// GitLab uses a dash as MR prefix for the (project-)internal ID (IID)
	return fmt.Sprintf("%s (!%d)", mergeRequest.Title, mergeRequest.IID)
}

func (c GitlabConfig) Driver(log logFn) (*GitlabDriver, error) {
	baseURL := gitlab.WithBaseURL(c.BaseURL())
	httpClient := gitlab.WithHTTPClient(&http.Client{})
	client, err := gitlab.NewOAuthClient(c.apiToken, httpClient, baseURL)
	if err != nil {
		return nil, err
	}
	driver := GitlabDriver{
		client:       client,
		GitlabConfig: c,
		log:          log,
	}
	return &driver, nil
}

func (c GitlabConfig) HostingServiceName() string {
	return "GitLab"
}

func (c GitlabConfig) NewProposalURL(branch, parentBranch string) (string, error) {
	query := url.Values{}
	query.Add("merge_request[source_branch]", branch)
	query.Add("merge_request[target_branch]", parentBranch)
	return fmt.Sprintf("%s/merge_requests/new?%s", c.RepositoryURL(), query.Encode()), nil
}

func (c GitlabConfig) ProjectPath() string {
	return fmt.Sprintf("%s/%s", c.owner, c.repository)
}

func (c GitlabConfig) RepositoryURL() string {
	return fmt.Sprintf("%s/%s", c.BaseURL(), c.ProjectPath())
}

// GitlabDriver provides access to the GitLab API.
type GitlabDriver struct {
	GitlabConfig
	client *gitlab.Client
	log    logFn
}

func (d *GitlabDriver) ProposalDetails(branch, parentBranch string) (*PullRequestInfo, error) {
	if d.apiToken == "" {
		return nil, nil //nolint:nilnil // we really want to return nil here
	}
	mergeRequests, err := d.loadMergeRequests(branch, parentBranch)
	if err != nil {
		return nil, err
	}
	if len(mergeRequests) < 1 {
		return nil, fmt.Errorf("no merge request from branch %q to branch %q found", branch, parentBranch)
	}
	if len(mergeRequests) > 1 {
		return nil, fmt.Errorf("found %d merge requests from branch %q to branch %q", len(mergeRequests), branch, parentBranch)
	}
	return &PullRequestInfo{
		CanMergeWithAPI:        true,
		DefaultProposalMessage: d.defaultCommitMessage(mergeRequests[0]),
		ProposalNumber:         mergeRequests[0].IID,
	}, nil
}

func (d *GitlabDriver) loadMergeRequests(branch, parentBranch string) ([]*gitlab.MergeRequest, error) {
	opts := &gitlab.ListProjectMergeRequestsOptions{
		State:        gitlab.String("opened"),
		SourceBranch: gitlab.String(branch),
		TargetBranch: gitlab.String(parentBranch),
	}
	// ListProjectMergeRequests takes care of encoding the project path already.
	mergeRequests, _, err := d.client.MergeRequests.ListProjectMergeRequests(d.ProjectPath(), opts)
	return mergeRequests, err
}

//nolint:nonamedreturns  // return value isn't obvious from function name
func (d *GitlabDriver) SquashMergeProposal(options SquashMergeProposalOptions) (mergeSha string, err error) {
	err = d.updatePullRequestsAgainst(options)
	if err != nil {
		return "", err
	}
	if options.ProposalNumber <= 0 {
		return "", fmt.Errorf("cannot merge via GitLab since there is no merge request")
	}
	if options.LogRequests {
		d.log("GitLab API: Merging MR !%d\n", options.ProposalNumber)
	}
	// GitLab API wants the full commit message in the body
	result, _, err := d.client.MergeRequests.AcceptMergeRequest(d.ProjectPath(), options.ProposalNumber, &gitlab.AcceptMergeRequestOptions{
		SquashCommitMessage: gitlab.String(options.CommitMessage),
		Squash:              gitlab.Bool(true),
		// This will be deleted by Git Town and make it fail if it is already deleted
		ShouldRemoveSourceBranch: gitlab.Bool(false),
		// SHA: gitlab.String(mergeSha),
	})
	if err != nil {
		return "", err
	}
	return result.SHA, nil
}

func (d *GitlabDriver) updatePullRequestsAgainst(options SquashMergeProposalOptions) error {
	// Fetch all open child merge requests that have this branch as their parent
	mergeRequests, _, err := d.client.MergeRequests.ListProjectMergeRequests(d.ProjectPath(), &gitlab.ListProjectMergeRequestsOptions{
		TargetBranch: gitlab.String(options.Branch),
		State:        gitlab.String("opened"),
	})
	if err != nil {
		return err
	}
	for _, mergeRequest := range mergeRequests {
		if options.LogRequests {
			d.log("GitLab API: Updating target branch for MR !%d\n", mergeRequest.IID)
		}
		// Update the target branch to be the latest version of the branch this MR is merged into
		_, _, err = d.client.MergeRequests.UpdateMergeRequest(d.ProjectPath(), mergeRequest.IID, &gitlab.UpdateMergeRequestOptions{
			TargetBranch: gitlab.String(options.ParentBranch),
		})
		if err != nil {
			return err
		}
	}
	return nil
}
