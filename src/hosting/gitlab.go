package hosting

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/git-town/git-town/v7/src/giturl"
	"github.com/xanzy/go-gitlab"
)

// GitlabDriver provides access to the API of GitLab installations.
type GitlabDriver struct {
	apiToken   string
	client     *gitlab.Client
	hostname   string
	originURL  string
	owner      string
	repository string
}

// NewGitlabDriver provides a GitLab driver instance if the given repo configuration is for a GitLab repo,
// otherwise nil.
func NewGitlabDriver(config config) *GitlabDriver {
	driverType := config.HostingService()
	originURL := config.OriginURL()
	hostname := giturl.Host(originURL)
	manualHostName := config.OriginOverride()
	if manualHostName != "" {
		hostname = manualHostName
	}
	if driverType != "gitlab" && hostname != "gitlab.com" {
		return nil
	}
	repositoryParts := strings.Split(giturl.Repo(originURL), "/")
	if len(repositoryParts) < 2 {
		return nil
	}
	lastIdx := len(repositoryParts) - 1
	owner := strings.Join(repositoryParts[:lastIdx], "/")
	repository := repositoryParts[lastIdx]
	return &GitlabDriver{
		apiToken:   config.GitLabToken(),
		originURL:  originURL,
		hostname:   hostname,
		owner:      owner,
		repository: repository,
	}
}

func (d *GitlabDriver) LoadPullRequestInfo(branch, parentBranch string) (result PullRequestInfo, err error) {
	if d.apiToken == "" {
		return result, nil
	}
	d.connect()
	mergeRequests, err := d.loadMergeRequests(branch, parentBranch)
	if err != nil {
		return result, err
	}
	if len(mergeRequests) != 1 {
		return result, err
	}
	result.CanMergeWithAPI = true
	result.DefaultCommitMessage = d.defaultCommitMessage(mergeRequests[0])
	result.PullRequestNumber = int64(mergeRequests[0].IID)
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

func (d *GitlabDriver) MergePullRequest(options MergePullRequestOptions) (mergeSha string, err error) {
	return "", errors.New("shipping pull requests via the GitLab API is currently not supported. If you need this functionality, please vote for it by opening a ticket at https://github.com/git-town/git-town/issues")
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
