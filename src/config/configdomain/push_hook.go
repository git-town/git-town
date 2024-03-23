package configdomain

import (
	"fmt"
	"strconv"

	"github.com/git-town/git-town/v13/src/gohacks"
	"github.com/git-town/git-town/v13/src/messages"
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

func NewPushHook(value, source string) (PushHook, error) {
	parsed, err := gohacks.ParseBool(value)
	if err != nil {
		return PushHook(true), fmt.Errorf(messages.ValueInvalid, source, value)
	}
	result := PushHook(parsed)
	return result, nil
}

func NewPushHookRef(value, source string) (*PushHook, error) {
	result, err := NewPushHook(value, source)
	return &result, err
}

// NoPushHook helps using the type checker to verify correct negation of the push-hook configuration setting.
type NoPushHook bool

func (noPushHook NoPushHook) Bool() bool {
	return bool(noPushHook)
}
