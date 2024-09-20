package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// BitbucketAppPassword is a bearer token to use with the GitHub API.
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
