package browserdomain

import "strconv"

type BrowserEnabled bool

func (self BrowserEnabled) Disabled() bool {
	return !self.Enabled()
}

func (self BrowserEnabled) Enabled() bool {
	return bool(self)
}

func (self BrowserEnabled) String() string {
	return strconv.FormatBool(bool(self))
}

func (self BrowserEnabled) StringHumanized() string {
	if self.Enabled() {
		return "enabled"
	}
	return "disabled"
}
