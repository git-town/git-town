package drivers

import (
	"fmt"
	"net/url"
	"strings"
)

var gitlabCodeHostingDriver = &CodeHostingDriver{

	CanBeUsed: func(hostname string) bool {
		return hostname == "gitlab.com" || strings.Contains(hostname, "gitlab")
	},

	GetNewPullRequestURL: func(repository string, branch string, parentBranch string) string {
		query := url.Values{}
		query.Add("merge_request[source_branch]", branch)
		query.Add("merge_request[target_branch]", parentBranch)
		return fmt.Sprintf("https://gitlab.com/%s/merge_requests/new?%s", repository, query.Encode())
	},

	GetRepositoryURL: func(repository string) string {
		return "https://gitlab.com/" + repository
	},

	HostingServiceName: "Gitlab",
}

func init() {
	registry.RegisterDriver(gitlabCodeHostingDriver)
}
