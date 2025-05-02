package configdomain

import (
	"strconv"

	"github.com/git-town/git-town/v20/internal/gohacks"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// PushHook contains the push-hook configuration setting.
type PushHook bool

func (self PushHook) IsTrue() bool {
	return bool(self)
}

func (self PushHook) Negate() NoPushHook {
	boolValue := bool(self)
	return NoPushHook(!boolValue)
}

func (self PushHook) String() string {
	return strconv.FormatBool(bool(self))
}

func ParsePushHook(value string, source Key) (Option[PushHook], error) {
	parsedOpt, err := gohacks.ParseBoolOpt(value, source.String())
	if parsed, has := parsedOpt.Get(); has {
		return Some(PushHook(parsed)), err
	}
	return None[PushHook](), err
}

// NoPushHook helps using the type checker to verify correct negation of the push-hook configuration setting.
type NoPushHook bool
