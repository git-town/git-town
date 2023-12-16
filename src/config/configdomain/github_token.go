package configdomain

// GitHubToken is a bearer token to use with the GitHub API.
type GitHubToken string

func (self GitHubToken) String() string {
	return string(self)
}

func NewGitHubTokenRef(value string) *GitHubToken {
	token := GitHubToken(value)
	return &token
}
