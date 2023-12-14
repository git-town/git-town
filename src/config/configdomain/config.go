package configdomain

import "github.com/git-town/git-town/v11/src/domain"

// Config is the merged configuration to be used by Git Town commands.
type Config struct {
	CodeHostingPlatformName string
	GiteaToken              GiteaToken
	GitHubToken             GitHubToken
	GitLabToken             GitLabToken
	MainBranch              domain.LocalBranchName
	Offline                 Offline
}

// Merges the given PartialConfig into this configuration object.
func (self *Config) Merge(other PartialConfig) {
	if other.CodeHostingPlatformName != nil {
		self.CodeHostingPlatformName = *other.CodeHostingPlatformName
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
	if other.MainBranch != nil {
		self.MainBranch = *other.MainBranch
	}
	if other.Offline != nil {
		self.Offline = *other.Offline
	}
}

// DefaultConfig provides the default configuration data to use when nothing is configured.
func DefaultConfig() Config {
	return Config{
		CodeHostingPlatformName: "",
		GiteaToken:              "",
		GitLabToken:             "",
		GitHubToken:             "",
		MainBranch:              domain.EmptyLocalBranchName(),
		Offline:                 false,
	}
}
