package drivers

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/git-town/git-town/src/drivers/helpers"
)

type gitlabConfig interface {
	GetCodeHostingDriverName() string
	GetRemoteOriginURL() string
	GetCodeHostingOriginHostname() string
}

type gitlabCodeHostingDriver struct {
	originURL  string
	hostname   string
	repository string
}

// LoadGitlab provides a GitLab driver instance if the given repo configuration is for a GitLab repo,
// otherwise nil.
func LoadGitlab(config gitlabConfig) CodeHostingDriver {
	driverType := config.GetCodeHostingDriverName()
	originURL := config.GetRemoteOriginURL()
	hostname := helpers.GetURLHostname(originURL)
	configuredHostName := config.GetCodeHostingOriginHostname()
	if configuredHostName != "" {
		hostname = configuredHostName
	}
	if driverType != "gitlab" && hostname != "gitlab.com" {
		return nil
	}
	return &gitlabCodeHostingDriver{
		originURL:  originURL,
		hostname:   hostname,
		repository: helpers.GetURLRepositoryName(originURL),
	}
}

func (d *gitlabCodeHostingDriver) CanMergeViaDriver(branch, parentBranch string) (canMerge bool, defaultCommitMessage string, pullRequestNumber int64, err error) {
	return false, "", 0, nil
}

func (d *gitlabCodeHostingDriver) NewPullRequestURL(branch, parentBranch string) string {
	query := url.Values{}
	query.Add("merge_request[source_branch]", branch)
	query.Add("merge_request[target_branch]", parentBranch)
	return fmt.Sprintf("%s/merge_requests/new?%s", d.GetRepositoryURL(), query.Encode())
}

func (d *gitlabCodeHostingDriver) GetRepositoryURL() string {
	return fmt.Sprintf("https://%s/%s", d.hostname, d.repository)
}

func (d *gitlabCodeHostingDriver) MergePullRequest(options MergePullRequestOptions) (mergeSha string, err error) {
	return "", errors.New("shipping pull requests via the GitLab API is currently not supported. If you need this functionality, please vote for it by opening a ticket at https://github.com/git-town/git-town/issues")
}

func (d *gitlabCodeHostingDriver) HostingServiceName() string {
	return "GitLab"
}
