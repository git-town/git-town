package configdomain

// Config is the merged configuration to be used by Git Town commands.
type Config struct {
	GitHubToken GitHubToken
}

// Merges the given PartialConfig into this configuration object.
func (self *Config) Merge(other PartialConfig) {
	if other.GitHubToken != nil {
		self.GitHubToken = *other.GitHubToken
	}
}

// DefaultConfig provides the default configuration data to use when nothing is configured.
func DefaultConfig() Config {
	return Config{
		GitHubToken: "",
	}
}
