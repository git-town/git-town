package domain

// PushHook contains the push-hook configuration setting.
type PushHook bool

// NoPushHook helps using the type checker to verify correct negation of the push-hook configuration setting.
type NoPushHook bool

func (self PushHook) Bool() bool {
	return bool(self)
}

func (self PushHook) Negate() NoPushHook {
	boolValue := bool(self)
	return NoPushHook(!boolValue)
}

func (self NoPushHook) Bool() bool {
	return bool(self)
}
