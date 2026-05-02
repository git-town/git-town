package browserdomain

import "strconv"

const NoBrowser = BrowserExecutable("(none)")

type BrowserEnabled bool

func (self BrowserEnabled) Enabled() bool {
	return bool(self)
}

func (self BrowserEnabled) String() string {
	return strconv.FormatBool(bool(self))
}
