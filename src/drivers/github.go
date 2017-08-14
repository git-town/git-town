package drivers

import (
	"fmt"
	"strings"

	"github.com/Originate/git-town/src/git"
)

type githubCodeHostingDriver struct {
}

func (d githubCodeHostingDriver) CanBeUsed(hostname string) bool {
	return hostname == "github.com" || strings.Contains(hostname, "github")
}

func (d githubCodeHostingDriver) GetNewPullRequestURL(repository string, branch string, parentBranch string) string {
	toCompare := branch
	if parentBranch != git.GetMainBranch() {
		toCompare = parentBranch + "..." + branch
	}
	return fmt.Sprintf("https://github.com/%s/compare/%s?expand=1", repository, toCompare)
}

func (d githubCodeHostingDriver) GetRepositoryURL(repository string) string {
	return "https://github.com/" + repository
}

func (d githubCodeHostingDriver) HostingServiceName() string {
	return "Github"
}

func init() {
	registry.RegisterDriver(githubCodeHostingDriver{})
}
