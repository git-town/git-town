package configdomain

// GitHubToken is a bearer token to use with the GitHub API.
type GitHubToken string

func (self GitHubToken) String() string {
	return string(self)
}

func NewGitHubToken(value string) GitHubToken {
	if len(value) == 0 {
		panic("received empty GitHub token")
	}
	return GitHubToken(value)
}
