package drivers

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/Originate/git-town/src/git"
)

type gitlabCodeHostingDriver struct {
	originURL  string
	hostname   string
	repository string
}

func (d *gitlabCodeHostingDriver) CanBeUsed(driverType string) bool {
	return driverType == "gitlab" || d.hostname == "gitlab.com"
}

func (d *gitlabCodeHostingDriver) CanMergePullRequest(branch, parentBranch string) (bool, string, error) {
	return false, "", nil
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

func (d *gitlabCodeHostingDriver) MergePullRequest(options MergePullRequestOptions) (string, error) {
	return "", errors.New("shipping pull requests via the Gitlab API is currently not supported. If you need this functionality, please vote for it by opening a ticket at https://github.com/originate/git-town/issues")
}

func (d *gitlabCodeHostingDriver) HostingServiceName() string {
	return "Gitlab"
}

func (d *gitlabCodeHostingDriver) SetOriginURL(originURL string) {
	d.originURL = originURL
	d.hostname = git.GetURLHostname(originURL)
	d.repository = git.GetURLRepositoryName(originURL)
}

func (d *gitlabCodeHostingDriver) SetOriginHostname(originHostname string) {
	d.hostname = originHostname
}

func (d *gitlabCodeHostingDriver) GetAPITokenKey() string {
	return ""
}

func (d *gitlabCodeHostingDriver) SetAPIToken(apiToken string) {}

func init() {
	registry.RegisterDriver(&gitlabCodeHostingDriver{})
}
