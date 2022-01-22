package hosting

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/git-town/git-town/v7/src/giturl"
)

// GitlabDriver provides access to the API of GitLab installations.
type GitlabDriver struct {
	hostname   string
	originURL  string
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
	return &GitlabDriver{
		originURL:  originURL,
		hostname:   hostname,
		repository: giturl.Repo(originURL),
	}
}

func (d *GitlabDriver) LoadPullRequestInfo(branch, parentBranch string) (result PullRequestInfo, err error) {
	return result, nil
}

func (d *GitlabDriver) NewPullRequestURL(branch, parentBranch string) (string, error) {
	query := url.Values{}
	query.Add("merge_request[source_branch]", branch)
	query.Add("merge_request[target_branch]", parentBranch)
	return fmt.Sprintf("%s/merge_requests/new?%s", d.RepositoryURL(), query.Encode()), nil
}

func (d *GitlabDriver) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s", d.hostname, d.repository)
}

func (d *GitlabDriver) MergePullRequest(options MergePullRequestOptions) (mergeSha string, err error) {
	return "", errors.New("shipping pull requests via the GitLab API is currently not supported. If you need this functionality, please vote for it by opening a ticket at https://github.com/git-town/git-town/issues")
}

func (d *GitlabDriver) HostingServiceName() string {
	return "GitLab"
}
