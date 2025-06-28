package forgedomain

import (
	"strings"

	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type BitbucketUsername string

func (self BitbucketUsername) String() string {
	return string(self)
}

func ParseBitbucketUsername(value string) Option[BitbucketUsername] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[BitbucketUsername]()
	}
	return Some(BitbucketUsername(value))
}
