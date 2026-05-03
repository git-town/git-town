package forgedomain

import (
	"strings"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// ForgejoToken is a bearer token to use with the Forgejo API.
type ForgejoToken string

func (self ForgejoToken) String() string {
	return string(self)
}

func ParseForgejoToken(valueOpt Option[string]) Option[ForgejoToken] {
	if value, has := valueOpt.Get(); has {
		return Some(ForgejoToken(strings.TrimSpace(value)))
	}
	return None[ForgejoToken]()
}
