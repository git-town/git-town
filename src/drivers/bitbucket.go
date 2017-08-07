package drivers

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Originate/git-town/src/git"
)

// BitbucketCodeHostingDriver provides functionality for working with
// repositories hosted on Bitbucket
var bitbucketCodeHostingDriver = &CodeHostingDriver{

	CanBeUsed: func(hostname string) bool {
		return hostname == "bitbucket.org" || strings.Contains(hostname, "bitbucket")
	},
	GetNewPullRequestURL: func(repository string, branch string, parentBranch string) string {
		query := url.Values{}
		query.Add("source", strings.Join([]string{repository, git.GetBranchSha(branch)[0:12], branch}, ":"))
		query.Add("dest", strings.Join([]string{repository, "", parentBranch}, ":"))
		return fmt.Sprintf("https://bitbucket.org/%s/pull-request/new?%s", repository, query.Encode())
	},

	GetRepositoryURL: func(repository string) string {
		return "https://bitbucket.org/" + repository
	},

	HostingServiceName: "Bitbucket",
}

func init() {
	registry.RegisterDriver(bitbucketCodeHostingDriver)
}
