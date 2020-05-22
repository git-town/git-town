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


func (d *giteaCodeHostingDriver) WasActivated(opts DriverOptions) bool {
	var hostname string

	if opts.OriginHostname != "" {
		hostname := opts.OriginHostname
	} else {
		hostname := git.Config().GetURLHostname(opts.OriginURL)
	}

	if opts.DriverType != "gitlab" && hostname != "gitlab.com"{
		return false
	}
    // Initialize
	d.hostname = hostname
	d.originURL = opts.OriginURL
	d.repository = git.Config().GetURLRepositoryName(opts.OriginURL)
	return true
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

func (d *gitlabCodeHostingDriver) GetAPIToken() string {
	return ""
}

func (d *gitlabCodeHostingDriver) SetAPIToken(apiToken string) {}

func init() {
	registry.RegisterDriver(&gitlabCodeHostingDriver{})
}
