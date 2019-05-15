package drivers

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/Originate/git-town/src/git"
)

type bitbucketCodeHostingDriver struct {
	originURL  string
	hostname   string
	repository string
}

func (d *bitbucketCodeHostingDriver) CanBeUsed(driverType string) bool {
	return driverType == "bitbucket" || d.hostname == "bitbucket.org"
}

func (d *bitbucketCodeHostingDriver) CanMergePullRequest(branch, parentBranch string) (bool, string, error) {
	return false, "", nil
}

func (d *bitbucketCodeHostingDriver) GetNewPullRequestURL(branch, parentBranch string) string {
	query := url.Values{}
	query.Add("source", strings.Join([]string{d.repository, git.GetBranchSha(branch)[0:12], branch}, ":"))
	query.Add("dest", strings.Join([]string{d.repository, "", parentBranch}, ":"))
	return fmt.Sprintf("%s/pull-request/new?%s", d.GetRepositoryURL(), query.Encode())
}

func (d *bitbucketCodeHostingDriver) GetRepositoryURL() string {
	return fmt.Sprintf("https://%s/%s", d.hostname, d.repository)
}

func (d *bitbucketCodeHostingDriver) MergePullRequest(options MergePullRequestOptions) (string, error) {
	return "", errors.New("shipping pull requests via the Bitbucket API is currently not supported. If you need this functionality, please vote for it by opening a ticket at https://github.com/originate/git-town/issues")
}

func (d *bitbucketCodeHostingDriver) HostingServiceName() string {
	return "Bitbucket"
}

func (d *bitbucketCodeHostingDriver) SetOriginURL(originURL string) {
	d.originURL = originURL
	d.hostname = git.GetURLHostname(originURL)
	d.repository = git.GetURLRepositoryName(originURL)
}

func (d *bitbucketCodeHostingDriver) SetOriginHostname(originHostname string) {
	d.hostname = originHostname
}

func (d *bitbucketCodeHostingDriver) GetAPITokenKey() string {
	return ""
}

func (d *bitbucketCodeHostingDriver) SetAPIToken(apiToken string) {}

func init() {
	registry.RegisterDriver(&bitbucketCodeHostingDriver{})
}
