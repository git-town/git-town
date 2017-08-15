package drivers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Originate/git-town/src/git"
)

type gitlabCodeHostingDriver struct {
	originURL string
	hostname  string
}

func (d *gitlabCodeHostingDriver) CanBeUsed() bool {
	return d.hostname == "gitlab.com" || strings.Contains(d.hostname, "gitlab")
}

func (d *gitlabCodeHostingDriver) GetNewPullRequestURL(repository string, branch string, parentBranch string) string {
	query := url.Values{}
	query.Add("merge_request[source_branch]", branch)
	query.Add("merge_request[target_branch]", parentBranch)
	return fmt.Sprintf("https://gitlab.com/%s/merge_requests/new?%s", repository, query.Encode())
}

func (d *gitlabCodeHostingDriver) GetRepositoryURL(repository string) string {
	return "https://gitlab.com/" + repository
}

func (d *gitlabCodeHostingDriver) HostingServiceName() string {
	return "Gitlab"
}

func (d *gitlabCodeHostingDriver) SetOriginURL(originURL string) {
	d.originURL = originURL
	d.hostname = git.GetURLHostname(originURL)
}

func init() {
	registry.RegisterDriver(&gitlabCodeHostingDriver{})
}
