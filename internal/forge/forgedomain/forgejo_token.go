package forgedomain

import (
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

// ForgejoToken is a bearer token to use with the Forgejo API.
type ForgejoToken string

func (self ForgejoToken) String() string {
	return string(self)
}

func ParseForgejoToken(value stringss.TrimmedString) Option[ForgejoToken] {
	if value == "" {
		return None[ForgejoToken]()
	}
	return Some(ForgejoToken(value))
}
