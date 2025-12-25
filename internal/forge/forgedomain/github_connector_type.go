package forgedomain

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// GitHubConnectorType describes the various ways in which Git Town can connect to the GitHub API.
type GitHubConnectorType string

const (
	GitHubConnectorTypeAPI GitHubConnectorType = "api" // connect to the GitHub API using Git Town's built-in API connector
	GitHubConnectorTypeGh  GitHubConnectorType = "gh"  // connect to the GitHub API by calling GitHub's "gh" tool
)

func (self GitHubConnectorType) String() string {
	return string(self)
}

// GitHubConnectorTypes provides all possible types that the GitHubConnectorTypes enum can have.
func GitHubConnectorTypes() []GitHubConnectorType {
	return []GitHubConnectorType{
		GitHubConnectorTypeAPI,
		GitHubConnectorTypeGh,
	}
}

func ParseGitHubConnectorType(text string, source string) (Option[GitHubConnectorType], error) {
	if text == "" {
		return None[GitHubConnectorType](), nil
	}
	for _, connectorType := range GitHubConnectorTypes() {
		if connectorType.String() == text {
			return Some(connectorType), nil
		}
	}
	return None[GitHubConnectorType](), fmt.Errorf(messages.GitHubConnectorTypeUnknown, source, text)
}
