package configdomain

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v14/src/gohacks"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
)

// SyncUpstream contains the configuration setting whether to sync with the upstream remote.
type SyncUpstream bool

func (self SyncUpstream) Bool() bool {
	return bool(self)
}

func (self SyncUpstream) String() string {
	return strconv.FormatBool(self.Bool())
}

func ParseSyncUpstream(value, source string) (SyncUpstream, error) {
	parsed, err := gohacks.ParseBool(value)
	if err != nil {
		return true, fmt.Errorf(messages.ValueInvalid, source, value)
	}
	return SyncUpstream(parsed), nil
}

func ParseSyncUpstreamOption(value, source string) (Option[SyncUpstream], error) {
	result, err := ParseSyncUpstream(value, source)
	if err != nil {
		return None[SyncUpstream](), err
	}
	return Some(result), err
}
