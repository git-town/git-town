package configdomain

import (
	"strconv"
)

// Offline is a new-type for the "offline" configuration setting.
type Offline bool

func (self Offline) IsOffline() bool {
	return bool(self)
}

func (self Offline) IsOnline() bool {
	return !self.IsOffline()
}

func (self Offline) String() string {
	return strconv.FormatBool(self.IsOffline())
}
