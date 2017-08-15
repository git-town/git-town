package drivers

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/Originate/git-town/src/git"
)

type gitlabCodeHostingDriver struct {
	originURL  string
	hostname   string
	repository string
}

func (d *gitlabCodeHostingDriver) CanBeUsed() bool {
	return d.hostname == "gitlab.com" || strings.Contains(d.hostname, "gitlab")
}

func (d *gitlabCodeHostingDriver) CanMergePullRequest(branch, parentBranch string) (bool, error) {
	return false, nil
}

func (d *gitlabCodeHostingDriver) GetNewPullRequestURL(branch, parentBranch string) string {
	query := url.Values{}
	query.Add("merge_request[source_branch]", branch)
	query.Add("merge_request[target_branch]", parentBranch)
	return fmt.Sprintf("%s/merge_requests/new?%s", d.GetRepositoryURL(), query.Encode())
}

func (d *gitlabCodeHostingDriver) GetRepositoryURL() string {
	return "https://gitlab.com/" + d.repository
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

func init() {
	registry.RegisterDriver(&gitlabCodeHostingDriver{})
}
