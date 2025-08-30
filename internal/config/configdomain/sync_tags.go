package configdomain

import "strconv"

// SyncTags contains the configuration setting whether to sync Git tags.
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
