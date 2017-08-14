package drivers

import (
	"fmt"
	"strings"

	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/util"
	"github.com/fatih/color"
)

// MergePullRequestOptions defines the options to the MergePullRequest function
type MergePullRequestOptions struct {
	Branch        string
	CommitMessage string
	LogRequests   bool
	ParentBranch  string
}

// CodeHostingDriver defines the interface
// of drivers for the different code hosting services
type CodeHostingDriver interface {
	CanMergePullRequest(branch, parentBranch string) (bool, error)
	GetRepositoryURL() string
	GetNewPullRequestURL(branch, parentBranch string) string
	MergePullRequest(options MergePullRequestOptions) (string, error)
}

// GetCodeHostingDriver returns an instance of the code hosting driver
// to use for the repository in the current working directory
func GetCodeHostingDriver() CodeHostingDriver {
	originURL := git.GetRemoteOriginURL()
	hostname := git.GetURLHostname(originURL)
	repository := git.GetURLRepositoryName(originURL)
	switch {
	case hostname == "github.com" || strings.Contains(hostname, "github"):
		return NewGithubCodeHostingDriver(repository)
	case hostname == "bitbucket.org" || strings.Contains(hostname, "bitbucket"):
		return NewBitbucketCodeHostingDriver(repository)
	case hostname == "gitlab.com" || strings.Contains(hostname, "gitlab"):
		return NewGitlabCodeHostingDriver(repository)
	default:
		util.ExitWithErrorMessage("Unsupported hosting service.", "This command requires hosting on GitHub, GitLab, or Bitbucket")
		return nil
	}
}

func printLog(message string) {
	fmt.Println()
	color.New(color.Bold).Println(message)
}
