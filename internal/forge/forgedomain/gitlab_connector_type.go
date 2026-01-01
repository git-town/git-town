package forgedomain

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// GitlabConnectorType describes the various ways in which Git Town can connect to the GitLab API.
type GitlabConnectorType string

const (
	GitLabConnectorTypeAPI  GitlabConnectorType = "api"  // connect to the GitLab API using Git Town's built-in API connector
	GitLabConnectorTypeGlab GitlabConnectorType = "glab" // connect to the GitLab API by calling GitLab's "glab" tool
)

func (self GitlabConnectorType) String() string {
	return string(self)
}

// GitLabConnectorTypes provides all possible types that the GitLabConnectorTypes enum can have.
func GitLabConnectorTypes() []GitlabConnectorType {
	return []GitlabConnectorType{
		GitLabConnectorTypeAPI,
		GitLabConnectorTypeGlab,
	}
}

func ParseGitLabConnectorType(text string, source string) (Option[GitlabConnectorType], error) {
	if text == "" {
		return None[GitlabConnectorType](), nil
	}
	for _, connectorType := range GitLabConnectorTypes() {
		if connectorType.String() == text {
			return Some(connectorType), nil
		}
	}
	return None[GitlabConnectorType](), fmt.Errorf(messages.GitLabConnectorTypeUnknown, source, text)
}
