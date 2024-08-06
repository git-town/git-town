package configdomain

import (
	"strconv"

	"github.com/git-town/git-town/v14/internal/gohacks"
	. "github.com/git-town/git-town/v14/pkg/prelude"
)

// PushHook contains the push-hook configuration setting.
type PushHook bool

func (pushHook PushHook) Bool() bool {
	return bool(pushHook)
}

func (pushHook PushHook) Negate() NoPushHook {
	boolValue := bool(pushHook)
	return NoPushHook(!boolValue)
}

func (pushHook PushHook) String() string {
	return strconv.FormatBool(pushHook.Bool())
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
