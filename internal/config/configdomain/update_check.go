package configdomain

import (
	"strconv"
)

// UpdateCheck is a new-type for the "update-check" configuration setting.
type UpdateCheck bool

func (self UpdateCheck) IsDisabled() bool {
	return !self.IsEnabled()
}

func (self UpdateCheck) IsEnabled() bool {
	return bool(self)
}

func (self UpdateCheck) String() string {
	return strconv.FormatBool(self.IsEnabled())
}
