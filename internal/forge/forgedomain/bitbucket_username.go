package forgedomain

import (
	"strings"

	. "github.com/git-town/git-town/v23/pkg/prelude"
)

type BitbucketUsername string

func (self BitbucketUsername) String() string {
	return string(self)
}

func ParseBitbucketUsername(valueOpt Option[string]) Option[BitbucketUsername] {
	if value, has := valueOpt.Get(); has {
		return Some(BitbucketUsername(strings.TrimSpace(value)))
	}
	return None[BitbucketUsername]()
}
