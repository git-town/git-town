package forgedomain

import (
	"strings"

	. "github.com/git-town/git-town/v23/pkg/prelude"
)

type BitbucketAppPassword string

func (self BitbucketAppPassword) String() string {
	return string(self)
}

func ParseBitbucketAppPassword(valueOpt Option[string]) Option[BitbucketAppPassword] {
	if value, has := valueOpt.Get(); has {
		return Some(BitbucketAppPassword(strings.TrimSpace(value)))
	}
	return None[BitbucketAppPassword]()
}
