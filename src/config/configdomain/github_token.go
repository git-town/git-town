package configdomain

import "github.com/git-town/git-town/v14/src/gohacks"

// GitHubToken is a bearer token to use with the GitHub API.
type GitHubToken string

func (self GitHubToken) String() string {
	return string(self)
}

func NewGitHubToken(value string) GitHubToken {
	if len(value) == 0 {
		panic("empty GitHub token")
	}
	return GitHubToken(value)
}

func NewGitHubTokenOption(value string) gohacks.Option[GitHubToken] {
	if value == "" {
		return gohacks.NewOptionNone[GitHubToken]()
	}
	return gohacks.NewOption(NewGitHubToken(value))
}
