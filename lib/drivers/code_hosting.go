package drivers

import (
	"strings"

	"github.com/Originate/git-town/lib/gitconfig"
	"github.com/Originate/git-town/lib/util"
)

type CodeHostingDriver interface {
	GetRepositoryUrl(repository string) string
	GetNewPullRequestUrl(repository string, branch string, parentBranch string) string
}

func GetCodeHostingDriver() CodeHostingDriver {
	hostname := gitconfig.GetUrlHostname(gitconfig.GetRemoteOriginUrl())
	if hostname == "github.com" || strings.Contains(hostname, "github") {
		return GithubCodeHostingDriver{}
	} else if hostname == "bitbucket.org" || strings.Contains(hostname, "bitbucket") {
		return BitbucketCodeHostingDriver{}
	} else if hostname == "gitlab.com" || strings.Contains(hostname, "gitlab") {
		return GitlabCodeHostingDriver{}
	} else {
		util.ExitWithErrorMessage("Unsupported hosting service.", "This command requires hosting on GitHub, GitLab, or Bitbucket")
		return nil
	}
}
