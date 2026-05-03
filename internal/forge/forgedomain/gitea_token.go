package forgedomain

import (
	"strings"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// GiteaToken is a bearer token to use with the Gitea API.
type GiteaToken string

func (self GiteaToken) String() string {
	return string(self)
}

func ParseGiteaToken(valueOpt Option[string]) Option[GiteaToken] {
	if value, has := valueOpt.Get(); has {
		return Some(GiteaToken(strings.TrimSpace(value)))
	}
	return None[GiteaToken]()
}
