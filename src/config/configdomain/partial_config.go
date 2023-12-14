package configdomain

// Data contains configuration data as it is stored in a particular configuration data source (Git, config file).
type PartialConfig struct {
	GitHubToken *GitHubToken
}

func (self *PartialConfig) Add(key Key, value string) bool {
	switch key {
	case KeyGithubToken:
		token := GitHubToken(value)
		self.GitHubToken = &token
		return true
	}
	return false
}

func EmptyPartialConfig() PartialConfig {
	return PartialConfig{}
}
