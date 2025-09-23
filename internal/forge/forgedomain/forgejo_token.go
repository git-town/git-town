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

func ParseForgejoToken(value string) Option[ForgejoToken] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[ForgejoToken]()
	}
	return Some(ForgejoToken(value))
}
