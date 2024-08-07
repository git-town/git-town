package configdomain

import (
	"strconv"

	"github.com/git-town/git-town/v15/internal/gohacks"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
)

// Offline is a new-type for the "offline" configuration setting.
type Offline bool

func (self Offline) Bool() bool {
	return bool(self)
}

func (self Offline) String() string {
	return strconv.FormatBool(self.Bool())
}

func (self Offline) ToOnline() Online {
	return Online(!self.Bool())
}

func ParseOffline(value, source string) (Option[Offline], error) {
	parsedOpt, err := gohacks.ParseBool(value, source)
	if parsed, has := parsedOpt.Get(); has {
		return Some(Offline(parsed)), err
	}
	return None[Offline](), err
}

type Online bool

func (online Online) Bool() bool {
	return bool(online)
}
