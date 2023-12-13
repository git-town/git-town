package domain

// PushHook contains the push-hook configuration setting.
type PushHook bool

func (self PushHook) Negate() PushHook {
	boolValue := bool(self)
	return PushHook(!boolValue)
}
