package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

// GiteaToken is a bearer token to use with the Gitea API.
type GiteaToken string

func (self GiteaToken) String() string {
	return string(self)
}

func NewGiteaToken(value string) GiteaToken {
	value = strings.TrimSpace(value)
	if value == "" {
		panic("empty Gitea token")
	}
	return GiteaToken(value)
}

func NewGiteaTokenOption(value string) Option[GiteaToken] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[GiteaToken]()
	}
	return Some(NewGiteaToken(value))
}
