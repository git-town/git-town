package drivers

import (
	"fmt"

	"github.com/Originate/git-town/lib/config"
)

type GithubCodeHostingDriver struct{}

func (driver GithubCodeHostingDriver) GetNewPullRequestUrl(repository string, branch string, parentBranch string) string {
	toCompare := branch
	if parentBranch != config.GetMainBranch() {
		toCompare = parentBranch + "..." + branch
	}
	return fmt.Sprintf("https://github.com/%s/compare/%s?expand=1", repository, toCompare)
}

func (driver GithubCodeHostingDriver) GetRepositoryUrl(repository string) string {
	return "https://github.com/" + repository
}
