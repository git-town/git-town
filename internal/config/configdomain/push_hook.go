package configdomain

import "strconv"

// PushHook contains the push-hook configuration setting.
type PushHook bool

func (self PushHook) Negate() NoPushHook {
	boolValue := bool(self)
	return NoPushHook(!boolValue)
}

func (self PushHook) ShouldRunPushHook() bool {
	return bool(self)
}

func (self PushHook) String() string {
	return strconv.FormatBool(bool(self))
}

// NoPushHook helps using the type checker to verify correct negation of the push-hook configuration setting.
type NoPushHook bool
