package forgedomain

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// GitLabConnectorType describes the various ways in which Git Town can connect to the GitLab API.
type GitLabConnectorType string

const (
	GitLabConnectorTypeAPI  GitLabConnectorType = "api"  // connect to the GitLab API using Git Town's built-in API connector
	GitLabConnectorTypeGlab GitLabConnectorType = "glab" // connect to the GitLab API by calling GitLab's "glab" tool
)

func (self GitLabConnectorType) String() string {
	return string(self)
}

// GitLabConnectorTypes provides all possible types that the GitLabConnectorTypes enum can have.
func GitLabConnectorTypes() []GitLabConnectorType {
	return []GitLabConnectorType{
		GitLabConnectorTypeAPI,
		GitLabConnectorTypeGlab,
	}
}

func ParseGitLabConnectorType(text string, source string) (Option[GitLabConnectorType], error) {
	if text == "" {
		return None[GitLabConnectorType](), nil
	}
	for _, connectorType := range GitLabConnectorTypes() {
		if connectorType.String() == text {
			return Some(connectorType), nil
		}
	}
	return None[GitLabConnectorType](), fmt.Errorf(messages.GitLabConnectorTypeUnknown, source, text)
}
