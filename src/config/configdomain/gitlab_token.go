package configdomain

import "github.com/git-town/git-town/v14/src/gohacks"

// GitLabToken is a bearer token to use with the GitLab API.
type GitLabToken string

func (self GitLabToken) String() string {
	return string(self)
}

func NewGitLabToken(value string) GitLabToken {
	return GitLabToken(value)
}

func NewGitLabTokenOption(value string) gohacks.Option[GitLabToken] {
	if value == "" {
		return gohacks.NewOptionNone[GitLabToken]()
	}
	return gohacks.NewOption(NewGitLabToken(value))
}
