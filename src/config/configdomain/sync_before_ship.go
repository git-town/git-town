package configdomain

import (
	"strconv"

	"github.com/git-town/git-town/v14/src/gohacks"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

type SyncBeforeShip bool

func (self SyncBeforeShip) Bool() bool {
	return bool(self)
}

func (self SyncBeforeShip) String() string {
	return strconv.FormatBool(self.Bool())
}

func ParseSyncBeforeShip(value, source string) (Option[SyncBeforeShip], error) {
	parsedOpt, err := gohacks.ParseBoolOpt(value, source)
	if parsed, has := parsedOpt.Get(); has {
		return Some(SyncBeforeShip(parsed)), err
	}
	return None[SyncBeforeShip](), err
}
