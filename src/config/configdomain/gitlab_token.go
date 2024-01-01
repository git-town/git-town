package configdomain

// GitLabToken is a bearer token to use with the GitLab API.
type GitLabToken string

func (self GitLabToken) String() string {
	return string(self)
}

func NewGitLabTokenRef(value string) *GitLabToken {
	token := GitLabToken(value)
	return &token
}
