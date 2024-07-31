package configdomain

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v14/src/gohacks"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/src/messages"
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

func ParsePushHookOption(valueStr, source string) (Option[PushHook], error) {
	if valueStr == "" {
		return None[PushHook](), nil
	}
	valueBool, err := gohacks.ParseBool(valueStr)
	if err != nil {
		return None[PushHook](), fmt.Errorf(messages.ValueInvalid, source, valueStr)
	}
	return Some(PushHook(valueBool)), nil
}

// NoPushHook helps using the type checker to verify correct negation of the push-hook configuration setting.
type NoPushHook bool

func (noPushHook NoPushHook) Bool() bool {
	return bool(noPushHook)
}
