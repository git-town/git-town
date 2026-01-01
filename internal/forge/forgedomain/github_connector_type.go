package forgedomain

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// GithubConnectorType describes the various ways in which Git Town can connect to the GitHub API.
type GithubConnectorType string

const (
	GithubConnectorTypeAPI GithubConnectorType = "api" // connect to the GitHub API using Git Town's built-in API connector
	GithubConnectorTypeGh  GithubConnectorType = "gh"  // connect to the GitHub API by calling GitHub's "gh" tool
)

func (self GithubConnectorType) String() string {
	return string(self)
}

// GithubConnectorTypes provides all possible types that the GithubConnectorTypes enum can have.
func GithubConnectorTypes() []GithubConnectorType {
	return []GithubConnectorType{
		GithubConnectorTypeAPI,
		GithubConnectorTypeGh,
	}
}

func ParseGithubConnectorType(text string, source string) (Option[GithubConnectorType], error) {
	if text == "" {
		return None[GithubConnectorType](), nil
	}
	for _, connectorType := range GithubConnectorTypes() {
		if connectorType.String() == text {
			return Some(connectorType), nil
		}
	}
	return None[GithubConnectorType](), fmt.Errorf(messages.GithubConnectorTypeUnknown, source, text)
}
