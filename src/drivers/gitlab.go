package drivers

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/git-town/git-town/src/git"
)

type gitlabCodeHostingDriver struct {
	originURL  string
	hostname   string
	repository string
}

func (d *gitlabCodeHostingDriver) CanBeUsed(driverType string) bool {
	return driverType == "gitlab" || d.hostname == "gitlab.com"
}

func (d *gitlabCodeHostingDriver) CanMergePullRequest(branch, parentBranch string) (canMerge bool, defaultCommitMessage string, pullRequestNumber int64, err error) {
	return false, "", 0, nil
}

func (d *gitlabCodeHostingDriver) GetNewPullRequestURL(branch, parentBranch string) string {
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

func (d *gitlabCodeHostingDriver) SetOriginURL(originURL string) {
	d.originURL = originURL
	d.hostname = git.Config().GetURLHostname(originURL)
	d.repository = git.Config().GetURLRepositoryName(originURL)
}

func (d *gitlabCodeHostingDriver) SetOriginHostname(originHostname string) {
	d.hostname = originHostname
}

func (d *gitlabCodeHostingDriver) GetAPIToken() string {
	return ""
}

func (d *gitlabCodeHostingDriver) SetAPIToken(apiToken string) {}

func init() {
	registry.RegisterDriver(&gitlabCodeHostingDriver{})
}
