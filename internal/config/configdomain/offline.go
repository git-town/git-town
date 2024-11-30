package configdomain

import (
	"strconv"

	"github.com/git-town/git-town/v16/internal/gohacks"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// Offline is a new-type for the "offline" configuration setting.
type Offline bool

func (self Offline) IsFalse() bool {
	return !self.IsTrue()
}

func (self Offline) IsTrue() bool {
	return bool(self)
}

func (self Offline) String() string {
	return strconv.FormatBool(self.IsTrue())
}

func (self Offline) ToOnline() Online {
	return Online(!self.IsTrue())
}

func ParseOffline(value string, source Key) (Option[Offline], error) {
	parsedOpt, err := gohacks.ParseBool(value, source.String())
	if parsed, has := parsedOpt.Get(); has {
		return Some(Offline(parsed)), err
	}
	return None[Offline](), err
}

type Online bool

func (online Online) IsTrue() bool {
	return bool(online)
}
