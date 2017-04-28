package drivers

import (
	"fmt"

	"github.com/Originate/git-town/src/git"
)

// GithubCodeHostingDriver provides tools for working with repositories
// hosted on Github
type GithubCodeHostingDriver struct{}

// GetNewPullRequestURL returns the URL of the page
// to create a new pull request on Github
func (driver GithubCodeHostingDriver) GetNewPullRequestURL(repository string, branch string, parentBranch string) string {
	toCompare := branch
	if parentBranch != git.GetMainBranch() {
		toCompare = parentBranch + "..." + branch
	}
	return fmt.Sprintf("https://github.com/%s/compare/%s?expand=1", repository, toCompare)
}

// GetRepositoryURL returns the URL of the given repository on github.com
func (driver GithubCodeHostingDriver) GetRepositoryURL(repository string) string {
	return "https://github.com/" + repository
}
