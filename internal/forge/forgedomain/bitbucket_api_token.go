package forgedomain

import (
	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

type BitbucketAPIToken stringss.Trimmed

func (self BitbucketAPIToken) String() string {
	return string(self)
}

func ParseBitbucketAPIToken(value stringss.Trimmed) Option[BitbucketAPIToken] {
	if value == "" {
		return None[BitbucketAPIToken]()
	}
	return Some(BitbucketAPIToken(value))
}
