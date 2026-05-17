package forgedomain

import (
	"fmt"

	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	"github.com/git-town/git-town/v23/internal/messages"
	. "github.com/git-town/git-town/v23/pkg/prelude"
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

func ParseGitlabConnectorType(text stringss.Trimmed, source string) (Option[GitlabConnectorType], error) {
	if text == "" {
		return None[GitlabConnectorType](), nil
	}
	for _, connectorType := range GitlabConnectorTypes() {
		if connectorType.String() == text.String() {
			return Some(connectorType), nil
		}
	}
	return None[GitlabConnectorType](), fmt.Errorf(messages.GitlabConnectorTypeUnknown, source, text)
}

func ParseGitlabConnectorTypeOpt(valueOpt Option[string], source string) (Option[GitlabConnectorType], error) {
	if value, has := valueOpt.Get(); has {
		return ParseGitlabConnectorType(value, source)
	}
	return None[GitlabConnectorType](), nil
}
