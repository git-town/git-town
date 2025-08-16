package forgedomain

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// GitHubTokenType describes the various ways in which the user can provide the GitHub API token to Git Town.
type GitHubTokenType string

const (
	GitHubTokenTypeEnter  GitHubTokenType = "enter"  // user enters the token manually
	GitHubTokenTypeScript GitHubTokenType = "script" // user provides a script that Git Town calls to get the GitHubToken
)

func (self GitHubTokenType) String() string {
	return string(self)
}

// GitHubConnectorTypes provides all possible types that the GitHubConnectorTypes enum can have.
func GitHubTokenTypes() []GitHubTokenType {
	return []GitHubTokenType{
		GitHubTokenTypeEnter,
		GitHubTokenTypeScript,
	}
}

func ParseGitHubTokenType(text string) (Option[GitHubTokenType], error) {
	if text == "" {
		return None[GitHubTokenType](), nil
	}
	for _, tokenType := range GitHubTokenTypes() {
		if tokenType.String() == text {
			return Some(tokenType), nil
		}
	}
	return None[GitHubTokenType](), fmt.Errorf(messages.GitHubTokenTypeUnknown, text)
}
