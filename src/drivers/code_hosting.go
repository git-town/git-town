package drivers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Originate/git-town/src/git"
	"github.com/fatih/color"
)

// MergePullRequestOptions defines the options to the MergePullRequest function
type MergePullRequestOptions struct {
	Branch        string
	CommitMessage string
	LogRequests   bool
	Owner         string
	ParentBranch  string
	Repository    string
}

// CodeHostingDriver defines the interface
// of drivers for the different code hosting services
type CodeHostingDriver interface {
	CanMergePullRequest() bool
	GetRepositoryURL(repository string) string
	GetNewPullRequestURL(repository string, branch string, parentBranch string) string
	MergePullRequest(options MergePullRequestOptions) (string, error)
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
		return nil
	}
}

// ValidateHasCodeHostingDriver returns an error if there is no code hosting driver
func ValidateHasCodeHostingDriver() error {
	driver := GetCodeHostingDriver()
	if driver == nil {
		return errors.New("Unsupported hosting service. This command requires hosting on GitHub, GitLab, or Bitbucket")
	}
	return nil
}

func printLog(message string) {
	fmt.Println()
	color.New(color.Bold).Println(message)
}
