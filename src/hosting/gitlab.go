package hosting

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/git-town/git-town/v7/src/hosting/helpers"
)

// GitlabCodeHostingDriver provides access to the API of GitLab installations.
type GitlabCodeHostingDriver struct {
	originURL  string
	hostname   string
	repository string
}

// LoadGitlab provides a GitLab driver instance if the given repo configuration is for a GitLab repo,
// otherwise nil.
func LoadGitlab(config config) *GitlabCodeHostingDriver {
	driverType := config.DriverName()
	originURL := config.RemoteOriginURL()
	hostname := helpers.URLHostname(originURL)
	manualHostName := config.OriginHost()
	if manualHostName != "" {
		hostname = manualHostName
	}
	if driverType != "gitlab" && hostname != "gitlab.com" {
		return nil
	}
	return &GitlabCodeHostingDriver{
		originURL:  originURL,
		hostname:   hostname,
		repository: helpers.URLRepositoryName(originURL),
	}
}

func (d *GitlabCodeHostingDriver) LoadPullRequestInfo(branch, parentBranch string) (result PullRequestInfo, err error) {
	return result, nil
}

func (d *GitlabCodeHostingDriver) NewPullRequestURL(branch, parentBranch string) (string, error) {
	query := url.Values{}
	query.Add("merge_request[source_branch]", branch)
	query.Add("merge_request[target_branch]", parentBranch)
	return fmt.Sprintf("%s/merge_requests/new?%s", d.RepositoryURL(), query.Encode()), nil
}

func (d *GitlabCodeHostingDriver) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s", d.hostname, d.repository)
}

func (d *GitlabCodeHostingDriver) MergePullRequest(options MergePullRequestOptions) (mergeSha string, err error) {
	return "", errors.New("shipping pull requests via the GitLab API is currently not supported. If you need this functionality, please vote for it by opening a ticket at https://github.com/git-town/git-town/issues")
}

func (d *GitlabCodeHostingDriver) HostingServiceName() string {
	return "GitLab"
}
