package configdomain

import (
	"strconv"
)

// indicates that the user does not want to sync tags
type SyncTags bool

func (self SyncTags) IsFalse() bool {
	return !self.IsTrue()
}

func (self SyncTags) IsTrue() bool {
	return bool(self)
}

func (self SyncTags) String() string {
	return strconv.FormatBool(self.IsTrue())
}
