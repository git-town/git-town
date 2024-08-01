package configdomain

import (
	"strconv"

	"github.com/git-town/git-town/v14/src/gohacks"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

// Offline is a new-type for the "offline" configuration setting.
type Offline bool

func (offline Offline) Bool() bool {
	return bool(offline)
}

func (offline Offline) String() string {
	return strconv.FormatBool(offline.Bool())
}

func (offline Offline) ToOnline() Online {
	return Online(!offline.Bool())
}

func NewOfflineOption(value, source string) (Option[Offline], error) {
	parsedOpt, err := gohacks.ParseBoolOpt(value, source)
	if parsed, has := parsedOpt.Get(); has {
		return Some(Offline(parsed)), err
	}
	return None[Offline](), err
}

type Online bool

func (online Online) Bool() bool {
	return bool(online)
}
