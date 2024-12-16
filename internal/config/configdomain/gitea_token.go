package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v17/pkg/prelude"
)

// GiteaToken is a bearer token to use with the Gitea API.
type GiteaToken string

func (self GiteaToken) String() string {
	return string(self)
}

func ParseGiteaToken(value string) Option[GiteaToken] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[GiteaToken]()
	}
	return Some(GiteaToken(value))
}
