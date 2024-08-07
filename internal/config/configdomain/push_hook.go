package configdomain

import (
	"strconv"

	"github.com/git-town/git-town/v15/internal/gohacks"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
)

// PushHook contains the push-hook configuration setting.
type PushHook bool

func (self PushHook) Bool() bool {
	return bool(self)
}

func (self PushHook) Negate() NoPushHook {
	boolValue := bool(self)
	return NoPushHook(!boolValue)
}

func (self PushHook) String() string {
	return strconv.FormatBool(self.Bool())
}

func ParsePushHook(value, source string) (Option[PushHook], error) {
	parsedOpt, err := gohacks.ParseBool(value, source)
	if parsed, has := parsedOpt.Get(); has {
		return Some(PushHook(parsed)), err
	}
	return None[PushHook](), err
}

// NoPushHook helps using the type checker to verify correct negation of the push-hook configuration setting.
type NoPushHook bool

func (noPushHook NoPushHook) Bool() bool {
	return bool(noPushHook)
}
