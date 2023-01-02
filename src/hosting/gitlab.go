package hosting

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/git-town/git-town/v7/src/giturl"
	"github.com/xanzy/go-gitlab"
)

// GitlabDriver provides access to the API of GitLab installations.
type GitlabDriver struct {
	apiToken   string
	client     *gitlab.Client
	hostname   string
	log        logFn
	originURL  string
	owner      string
	repository string
}

// NewGitlabDriver provides a GitLab driver instance if the given repo configuration is for a GitLab repo,
// otherwise nil.
func NewGitlabDriver(url giturl.Parts, config config, hostingConfig hostingConfig, log logFn) *GitlabDriver {
	driverType := hostingConfig.HostingService()
	manualHostName := config.OriginOverride()
	if manualHostName != "" {
		url.Host = manualHostName
	}
	if driverType != "gitlab" && url.Host != "gitlab.com" {
		return nil
	}
	return &GitlabDriver{
		apiToken:   hostingConfig.GitLabToken(),
		originURL:  config.OriginURL(),
		hostname:   url.Host,
		log:        log,
		owner:      url.Org,
		repository: url.Repo,
	}
}

func (d *GitlabDriver) LoadPullRequestInfo(branch, parentBranch string) (PullRequestInfo, error) {
	if d.apiToken == "" {
		return PullRequestInfo{}, nil
	}
	d.connect()
	mergeRequests, err := d.loadMergeRequests(branch, parentBranch)
	if err != nil {
		return PullRequestInfo{}, err
	}
	if len(mergeRequests) != 1 {
		return PullRequestInfo{}, nil
	}
	result := PullRequestInfo{
		CanMergeWithAPI:      true,
		DefaultCommitMessage: d.defaultCommitMessage(mergeRequests[0]),
		PullRequestNumber:    int64(mergeRequests[0].IID),
	}
	return result, nil
}

func (d *GitlabDriver) NewPullRequestURL(branch, parentBranch string) (string, error) {
	query := url.Values{}
	query.Add("merge_request[source_branch]", branch)
	query.Add("merge_request[target_branch]", parentBranch)
	return fmt.Sprintf("%s/merge_requests/new?%s", d.RepositoryURL(), query.Encode()), nil
}

func (d *GitlabDriver) BaseURL() string {
	return fmt.Sprintf("https://%s", d.hostname)
}

func (d *GitlabDriver) ProjectPath() string {
	return fmt.Sprintf("%s/%s", d.owner, d.repository)
}

func (d *GitlabDriver) RepositoryURL() string {
	return fmt.Sprintf("%s/%s", d.BaseURL(), d.ProjectPath())
}

//nolint:nonamedreturns  // return value isn't obvious from function name
func (d *GitlabDriver) MergePullRequest(options MergePullRequestOptions) (mergeSha string, err error) {
	d.connect()
	err = d.updatePullRequestsAgainst(options)
	if err != nil {
		return "", err
	}
	return d.mergePullRequest(options)
}

func (d *GitlabDriver) HostingServiceName() string {
	return "GitLab"
}

// Helper

func (d *GitlabDriver) connect() {
	if d.client == nil {
		baseURL := gitlab.WithBaseURL(d.BaseURL())
		httpClient := gitlab.WithHTTPClient(&http.Client{})
		client, err := gitlab.NewOAuthClient(d.apiToken, httpClient, baseURL)
		if err == nil {
			d.client = client
		}
	}
}

func (d *GitlabDriver) defaultCommitMessage(mergeRequest *gitlab.MergeRequest) string {
	// GitLab uses a dash as MR prefix for the (project-)internal ID (IID)
	return fmt.Sprintf("%s (!%d)", mergeRequest.Title, mergeRequest.IID)
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
func (d *GitlabDriver) mergePullRequest(options MergePullRequestOptions) (mergeSha string, err error) {
	if options.PullRequestNumber <= 0 {
		return "", fmt.Errorf("cannot merge via GitLab since there is no merge request")
	}
	if options.LogRequests {
		d.log("GitLab API: Merging MR !%d\n", options.PullRequestNumber)
	}
	// GitLab API wants the full commit message in the body
	result, _, err := d.client.MergeRequests.AcceptMergeRequest(d.ProjectPath(), int(options.PullRequestNumber), &gitlab.AcceptMergeRequestOptions{
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

func (d *GitlabDriver) updatePullRequestsAgainst(options MergePullRequestOptions) error {
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
