package drivers

import (
	"fmt"
	"net/url"
	"strings"
)

type gitlabCodeHostingDriver struct {
}

func (d gitlabCodeHostingDriver) CanBeUsed(hostname string) bool {
	return hostname == "gitlab.com" || strings.Contains(hostname, "gitlab")
}

func (d gitlabCodeHostingDriver) GetNewPullRequestURL(repository string, branch string, parentBranch string) string {
	query := url.Values{}
	query.Add("merge_request[source_branch]", branch)
	query.Add("merge_request[target_branch]", parentBranch)
	return fmt.Sprintf("https://gitlab.com/%s/merge_requests/new?%s", repository, query.Encode())
}

func (d gitlabCodeHostingDriver) GetRepositoryURL(repository string) string {
	return "https://gitlab.com/" + repository
}

func (d gitlabCodeHostingDriver) HostingServiceName() string {
	return "Gitlab"
}

func init() {
	registry.RegisterDriver(gitlabCodeHostingDriver{})
}
