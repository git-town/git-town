package drivers

import (
	"fmt"
	"strings"

	"github.com/Originate/git-town/src/git"
)

var githubCodeHostingDriver = &CodeHostingDriver{

	CanBeUsed: func(hostname string) bool {
		return hostname == "github.com" || strings.Contains(hostname, "github")
	},

	GetNewPullRequestURL: func(repository string, branch string, parentBranch string) string {
		toCompare := branch
		if parentBranch != git.GetMainBranch() {
			toCompare = parentBranch + "..." + branch
		}
		return fmt.Sprintf("https://github.com/%s/compare/%s?expand=1", repository, toCompare)
	},

	GetRepositoryURL: func(repository string) string {
		return "https://github.com/" + repository
	},

	HostingServiceName: "Github",
}

func init() {
	registry.RegisterDriver(githubCodeHostingDriver)
}
