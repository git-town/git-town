package configdomain

import "strconv"

// SyncUpstream contains the configuration setting whether to sync with the upstream remote.
type SyncUpstream bool

func (self SyncUpstream) ShouldSyncUpstream() bool {
	return bool(self)
}

func (self SyncUpstream) String() string {
	return strconv.FormatBool(bool(self))
}
