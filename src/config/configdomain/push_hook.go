package configdomain

import "strconv"

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

// NoPushHook helps using the type checker to verify correct negation of the push-hook configuration setting.
type NoPushHook bool

func (noPushHook NoPushHook) Bool() bool {
	return bool(noPushHook)
}
