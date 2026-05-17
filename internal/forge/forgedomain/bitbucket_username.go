package forgedomain

import (
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

type BitbucketUsername string

func (self BitbucketUsername) String() string {
	return string(self)
}

func ParseBitbucketUsername(value stringss.TrimmedString) Option[BitbucketUsername] {
	if value == "" {
		return None[BitbucketUsername]()
	}
	return Some(BitbucketUsername(value))
}
