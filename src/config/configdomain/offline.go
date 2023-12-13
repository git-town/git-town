package configdomain

// Offline is a new-type for the "offline" configuration setting.
type Offline bool

func (self Offline) Bool() bool {
	return bool(self)
}
