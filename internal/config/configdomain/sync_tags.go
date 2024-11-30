package configdomain

import (
	"strconv"

	"github.com/git-town/git-town/v16/internal/gohacks"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

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

func ParseSyncTags(value string, source Key) (Option[SyncTags], error) {
	parsedOpt, err := gohacks.ParseBool(value, source.String())
	if parsed, has := parsedOpt.Get(); has {
		return Some(SyncTags(parsed)), err
	}
	return None[SyncTags](), err
}
