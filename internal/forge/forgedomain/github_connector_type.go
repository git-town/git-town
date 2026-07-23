package forgedomain

import (
	"fmt"

	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	"github.com/git-town/git-town/v24/internal/messages"
	. "github.com/git-town/git-town/v24/pkg/prelude"
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

func ParseGithubConnectorType(text stringss.Trimmed, source string) (Option[GithubConnectorType], error) {
	if text == "" {
		return None[GithubConnectorType](), nil
	}
	for _, connectorType := range GithubConnectorTypes() {
		if connectorType.String() == text.String() {
			return Some(connectorType), nil
		}
	}
	return None[GithubConnectorType](), fmt.Errorf(messages.GithubConnectorTypeUnknown, source, text)
}

func ParseGithubConnectorTypeOpt(valueOpt Option[string], source string) (Option[GithubConnectorType], error) {
	value, has := valueOpt.Get()
	if !has {
		return None[GithubConnectorType](), nil
	}
	return ParseGithubConnectorType(value, source)
}
