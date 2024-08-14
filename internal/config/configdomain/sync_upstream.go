package configdomain

import (
	"strconv"

	"github.com/git-town/git-town/v15/internal/gohacks"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
)

// SyncUpstream contains the configuration setting whether to sync with the upstream remote.
type SyncUpstream bool

func (self SyncUpstream) IsTrue() bool {
	return bool(self)
}

func (self SyncUpstream) String() string {
	return strconv.FormatBool(bool(self))
}

func ParseSyncUpstream(value, source string) (Option[SyncUpstream], error) {
	parsedOpt, err := gohacks.ParseBool(value, source)
	if parsed, has := parsedOpt.Get(); has {
		return Some(SyncUpstream(parsed)), err
	}
	return None[SyncUpstream](), err
}
