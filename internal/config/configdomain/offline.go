package configdomain

import (
	"strconv"

	"github.com/git-town/git-town/v21/internal/gohacks"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Offline is a new-type for the "offline" configuration setting.
type Offline bool

func (self Offline) IsOffline() bool {
	return bool(self)
}

func (self Offline) IsOnline() bool {
	return !self.IsOffline()
}

func (self Offline) String() string {
	return strconv.FormatBool(self.IsOffline())
}

func ParseOffline(value string, source Key) (Option[Offline], error) {
	parsedOpt, err := gohacks.ParseBoolOpt(value, source.String())
	if parsed, has := parsedOpt.Get(); has {
		return Some(Offline(parsed)), err
	}
	return None[Offline](), err
}
