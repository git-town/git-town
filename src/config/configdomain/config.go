package configdomain

// Config is the merged configuration to be used by Git Town commands.
type Config struct {
	CodeHostingPlatformName *string
	GiteaToken              GiteaToken
	GitHubToken             GitHubToken
	GitLabToken             GitLabToken
}

// Merges the given PartialConfig into this configuration object.
func (self *Config) Merge(other PartialConfig) {
	if other.CodeHostingPlatformName != nil {
		self.CodeHostingPlatformName = other.CodeHostingPlatformName
	}
	if other.GiteaToken != nil {
		self.GiteaToken = *other.GiteaToken
	}
	if other.GitHubToken != nil {
		self.GitHubToken = *other.GitHubToken
	}
	if other.GitLabToken != nil {
		self.GitLabToken = *other.GitLabToken
	}
}

// DefaultConfig provides the default configuration data to use when nothing is configured.
func DefaultConfig() Config {
	emptyString := ""
	return Config{ //nolint:exhaustruct
		CodeHostingPlatformName: &emptyString,
		GiteaToken:              "",
		GitLabToken:             "",
		GitHubToken:             "",
	}
}
