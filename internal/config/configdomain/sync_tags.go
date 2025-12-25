package configdomain

import "strconv"

// SyncTags contains the configuration setting whether to sync Git tags.
type SyncTags bool

func (self SyncTags) ShouldSyncTags() bool {
	return bool(self)
}

func (self SyncTags) String() string {
	return strconv.FormatBool(self.ShouldSyncTags())
}
