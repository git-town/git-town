package configdomain

import (
	"strconv"

	"github.com/git-town/git-town/v15/internal/gohacks"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
)

// indicates that the user does not want to sync tags
type SyncTags bool

func (self SyncTags) Bool() bool {
	return bool(self)
}

func (self SyncTags) String() string {
	return strconv.FormatBool(self.Bool())
}

func (self SyncTags) ToOnline() Online {
	return Online(!self.Bool())
}

func ParseNoTags(value, source string) (Option[SyncTags], error) {
	parsedOpt, err := gohacks.ParseBool(value, source)
	if parsed, has := parsedOpt.Get(); has {
		return Some(SyncTags(parsed)), err
	}
	return None[SyncTags](), err
}
