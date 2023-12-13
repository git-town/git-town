package configdomain

// Data contains configuration data as it is stored in a particular configuration data source (Git, config file).
type PartialConfig struct {
	GitHubToken *GitHubToken
}
