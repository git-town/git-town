package forgedomain

import (
	"strings"

	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// CodebergToken is a bearer token to use with the Codeberg API.
type CodebergToken string

func (self CodebergToken) String() string {
	return string(self)
}

func ParseCodebergToken(value string) Option[CodebergToken] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[CodebergToken]()
	}
	return Some(CodebergToken(value))
}
