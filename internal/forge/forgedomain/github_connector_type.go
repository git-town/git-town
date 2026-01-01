package forgedomain

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// GithubConnectorType describes the various ways in which Git Town can connect to the GitHub API.
type GithubConnectorType string

const (
	GitHubConnectorTypeAPI GithubConnectorType = "api" // connect to the GitHub API using Git Town's built-in API connector
	GitHubConnectorTypeGh  GithubConnectorType = "gh"  // connect to the GitHub API by calling GitHub's "gh" tool
)

func (self GithubConnectorType) String() string {
	return string(self)
}

// GitHubConnectorTypes provides all possible types that the GitHubConnectorTypes enum can have.
func GitHubConnectorTypes() []GithubConnectorType {
	return []GithubConnectorType{
		GitHubConnectorTypeAPI,
		GitHubConnectorTypeGh,
	}
}

func ParseGitHubConnectorType(text string, source string) (Option[GithubConnectorType], error) {
	if text == "" {
		return None[GithubConnectorType](), nil
	}
	for _, connectorType := range GitHubConnectorTypes() {
		if connectorType.String() == text {
			return Some(connectorType), nil
		}
	}
	return None[GithubConnectorType](), fmt.Errorf(messages.GitHubConnectorTypeUnknown, source, text)
}
