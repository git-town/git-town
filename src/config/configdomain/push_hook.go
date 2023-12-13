package configdomain

// PushHook contains the push-hook configuration setting.
type PushHook bool

// NoPushHook helps using the type checker to verify correct negation of the push-hook configuration setting.
type NoPushHook bool

func (pushHook PushHook) Bool() bool {
	return bool(pushHook)
}

func (pushHook PushHook) Negate() NoPushHook {
	boolValue := bool(pushHook)
	return NoPushHook(!boolValue)
}

func (noPushHook NoPushHook) Bool() bool {
	return bool(noPushHook)
}
