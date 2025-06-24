package forgedomain

// GitHubConnectorType describes the various ways in which Git Town can connect to the GitHub API.
type GitHubConnectorType string

const (
	GitHubConnectorTypeAPI GitHubConnectorType = "api" // connect to the GitHub API using Git Town's built-in API connector
	GitHubConnectorTypeGh  GitHubConnectorType = "gh"  // connect to the GitHub API by calling GitHub's "gh" tool
)

func (self GitHubConnectorType) String() string {
	return string(self)
}
