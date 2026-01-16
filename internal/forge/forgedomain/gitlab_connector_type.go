package forgedomain

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// GitlabConnectorType describes the various ways in which Git Town can connect to the GitLab API.
type GitlabConnectorType string

const (
	GitlabConnectorTypeAPI  GitlabConnectorType = "api"  // connect to the GitLab API using Git Town's built-in API connector
	GitlabConnectorTypeGlab GitlabConnectorType = "glab" // connect to the GitLab API by calling GitLab's "glab" tool
)

func (self GitlabConnectorType) String() string {
	return string(self)
}

// GitlabConnectorTypes provides all possible types that the GitlabConnectorTypes enum can have.
func GitlabConnectorTypes() []GitlabConnectorType {
	return []GitlabConnectorType{
		GitlabConnectorTypeAPI,
		GitlabConnectorTypeGlab,
	}
}

func ParseGitlabConnectorType(text string, source string) (Option[GitlabConnectorType], error) {
	if text == "" {
		return None[GitlabConnectorType](), nil
	}
	for _, connectorType := range GitlabConnectorTypes() {
		if connectorType.String() == text {
			return Some(connectorType), nil
		}
	}
	return None[GitlabConnectorType](), fmt.Errorf(messages.GitlabConnectorTypeUnknown, source, text)
}
