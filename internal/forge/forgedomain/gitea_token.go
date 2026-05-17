package forgedomain

import (
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

// GiteaToken is a bearer token to use with the Gitea API.
type GiteaToken stringss.Trimmed

func (self GiteaToken) String() string {
	return string(self)
}

func ParseGiteaToken(value stringss.Trimmed) Option[GiteaToken] {
	if value == "" {
		return None[GiteaToken]()
	}
	return Some(GiteaToken(value))
}
