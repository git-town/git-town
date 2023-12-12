package github

type gitTownConfig interface {
	// GitHubToken provides the personal access token for GitHub stored in the Git configuration.
	GitHubToken() string
}
