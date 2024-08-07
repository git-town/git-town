package configdomain

import (
	"strconv"

	"github.com/git-town/git-town/v15/internal/gohacks"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
)

// indicates that the user does not want to sync tags
type SyncTags bool

func (self SyncTags) IsTrue() bool {
	return bool(self)
}

func (self SyncTags) String() string {
	return strconv.FormatBool(self.IsTrue())
}

func ParseSyncTags(value, source string) (Option[SyncTags], error) {
	parsedOpt, err := gohacks.ParseBool(value, source)
	if parsed, has := parsedOpt.Get(); has {
		return Some(SyncTags(parsed)), err
	}
	return None[SyncTags](), err
}
