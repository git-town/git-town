package forgedomain

import (
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

type BitbucketAppPassword stringss.Trimmed

func (self BitbucketAppPassword) String() string {
	return string(self)
}

func ParseBitbucketAppPassword(value stringss.Trimmed) Option[BitbucketAppPassword] {
	if value == "" {
		return None[BitbucketAppPassword]()
	}
	return Some(BitbucketAppPassword(value))
}
