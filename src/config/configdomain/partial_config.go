package configdomain

// Data contains configuration data as it is stored in a particular configuration data source (Git, config file).
type PartialConfig struct {
	CodeHostingPlatformName *string
	GiteaToken              *GiteaToken
	GitHubToken             *GitHubToken
	GitLabToken             *GitLabToken
}

func (self *PartialConfig) Add(key Key, value string) bool {
	switch key {
	case KeyGiteaToken:
		token := GiteaToken(value)
		self.GiteaToken = &token
	case KeyGithubToken:
		token := GitHubToken(value)
		self.GitHubToken = &token
	case KeyGitlabToken:
		token := GitLabToken(value)
		self.GitLabToken = &token
	case KeyCodeHostingPlatform:
		self.CodeHostingPlatformName = &value
	default:
		return false
	}
	return true
}

func EmptyPartialConfig() PartialConfig {
	return PartialConfig{} //nolint:exhaustruct
}
