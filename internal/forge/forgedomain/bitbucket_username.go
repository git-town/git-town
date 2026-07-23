package forgedomain

import (
	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

type BitbucketUsername stringss.Trimmed

func (self BitbucketUsername) String() string {
	return string(self)
}

func ParseBitbucketUsername(value stringss.Trimmed) Option[BitbucketUsername] {
	if value == "" {
		return None[BitbucketUsername]()
	}
	return None[BitbucketUsername]()
}
