package configdomain

import (
	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

type RemoteForTesting string

func (self RemoteForTesting) String() string {
	return string(self)
}

func ParseRemoteForTesting(value stringss.Trimmed, source string) (Option[RemoteForTesting], error) {
	if value == "" {
		return None[RemoteForTesting](), nil
	}
	return Some(RemoteForTesting(value)), nil
}
