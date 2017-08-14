package drivers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Originate/git-town/src/git"
)

type bitbucketCodeHostingDriver struct {
	originURL string
	hostname  string
}

func (d *bitbucketCodeHostingDriver) CanBeUsed() bool {
	return d.hostname == "bitbucket.org" || strings.Contains(d.hostname, "bitbucket")
}

func (d *bitbucketCodeHostingDriver) GetNewPullRequestURL(repository string, branch string, parentBranch string) string {
	query := url.Values{}
	query.Add("source", strings.Join([]string{repository, git.GetBranchSha(branch)[0:12], branch}, ":"))
	query.Add("dest", strings.Join([]string{repository, "", parentBranch}, ":"))
	return fmt.Sprintf("https://bitbucket.org/%s/pull-request/new?%s", repository, query.Encode())
}

func (d *bitbucketCodeHostingDriver) GetRepositoryURL(repository string) string {
	return "https://bitbucket.org/" + repository
}

func (d *bitbucketCodeHostingDriver) HostingServiceName() string {
	return "Bitbucket"
}

func (d *bitbucketCodeHostingDriver) SetOriginURL(originURL string) {
	d.originURL = originURL
	d.hostname = git.GetURLHostname(originURL)
}

func init() {
	registry.RegisterDriver(&bitbucketCodeHostingDriver{})
}
