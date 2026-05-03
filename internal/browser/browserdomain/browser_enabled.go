package browserdomain

import (
	"strconv"

	. "github.com/git-town/git-town/v22/pkg/prelude"
)

const NoBrowser = BrowserExecutable("(none)")

type BrowserEnabled bool

func (self BrowserEnabled) Enabled() bool {
	return bool(self)
}

func (self BrowserEnabled) String() string {
	return strconv.FormatBool(bool(self))
}

func NewBrowserEnabledFromTTY(tty HasTTY) Option[BrowserEnabled] {
	if tty {
		return None[BrowserEnabled]()
	}
	return Some(BrowserEnabled(false))
}
