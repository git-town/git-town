package domain

type PushHook bool

func (self PushHook) Negate() PushHook {
	boolValue := bool(self)
	return PushHook(!boolValue)
}
