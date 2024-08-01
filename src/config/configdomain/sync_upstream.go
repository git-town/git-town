package configdomain

import (
	"strconv"

	"github.com/git-town/git-town/v14/src/gohacks"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

// SyncUpstream contains the configuration setting whether to sync with the upstream remote.
type SyncUpstream bool

func (self SyncUpstream) Bool() bool {
	return bool(self)
}

func (self SyncUpstream) String() string {
	return strconv.FormatBool(self.Bool())
}

func ParseSyncUpstreamOption(value, source string) (Option[SyncUpstream], error) {
	parsedOpt, err := gohacks.ParseBoolOpt(value, source)
	if parsed, has := parsedOpt.Get(); has {
		return Some(SyncUpstream(parsed)), err
	}
	return None[SyncUpstream](), err
}
