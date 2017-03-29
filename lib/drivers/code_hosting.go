package drivers

import (
	"strings"

	"github.com/Originate/git-town/lib/config"
	"github.com/Originate/git-town/lib/util"
)

type CodeHostingDriver interface {
	GetRepositoryUrl(repository string) string
	GetNewPullRequestUrl(repository string, branch string, parentBranch string) string
}

func GetCodeHostingDriver() CodeHostingDriver {
	originHostname := config.GetRemoteOriginHostname()
	if originHostname == "github.com" {
		return GithubCodeHostingDriver{}
	} else if originHostname == "bitbucket.org" {
		return BitbucketCodeHostingDriver{}
	} else if originHostname == "gitlab.com" {
		return GitlabCodeHostingDriver{}
	} else if strings.Contains(originHostname, "github") {
		return GithubCodeHostingDriver{}
	} else if strings.Contains(originHostname, "bitbucket") {
		return BitbucketCodeHostingDriver{}
	} else if strings.Contains(originHostname, "gitlab") {
		return GitlabCodeHostingDriver{}
	} else {
		util.ExitWithErrorMessage("Unsupported hosting service.", "This command requires hosting on GitHub, GitLab, or Bitbucket")
		return nil
	}
}
