package configdomain

import "strconv"

// IgnoreUncommitted indicates whether to stash uncommitted changes when shipping.
type IgnoreUncommitted bool

func (self IgnoreUncommitted) ShouldIgnoreUncommitted() bool {
	return bool(self)
}

func (self IgnoreUncommitted) String() string {
	return strconv.FormatBool(bool(self))
}
