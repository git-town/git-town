package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v19/pkg/prelude"
)

type BitbucketAppPassword string

func (self BitbucketAppPassword) String() string {
	return string(self)
}

func ParseBitbucketAppPassword(value string) Option[BitbucketAppPassword] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[BitbucketAppPassword]()
	}
	return Some(BitbucketAppPassword(value))
}
