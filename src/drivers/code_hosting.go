package drivers

import (
	"strings"

	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/util"
)

// CodeHostingDriver defines the interface
// of drivers for the different code hosting services
type CodeHostingDriver interface {
	GetRepositoryURL(repository string) string
	GetNewPullRequestURL(repository string, branch string, parentBranch string) string
}

// GetCodeHostingDriver returns an instance of the code hosting driver
// to use for the repository in the current working directory
func GetCodeHostingDriver() CodeHostingDriver {
	hostname := git.GetURLHostname(git.GetRemoteOriginURL())
	switch {
	case hostname == "github.com" || strings.Contains(hostname, "github"):
		return GithubCodeHostingDriver{}
	case hostname == "bitbucket.org" || strings.Contains(hostname, "bitbucket"):
		return BitbucketCodeHostingDriver{}
	case hostname == "gitlab.com" || strings.Contains(hostname, "gitlab"):
		return GitlabCodeHostingDriver{}
	default:
		util.ExitWithErrorMessage("Unsupported hosting service.", "This command requires hosting on GitHub, GitLab, or Bitbucket")
		return nil
	}
}
