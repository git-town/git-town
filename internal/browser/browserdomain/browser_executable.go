package browserdomain

import (
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

const NoBrowser = BrowserExecutable("(none)")

// BrowserExecutable indicates a custom browser to use.
// If set to "" or "(none)", browsers are disabled.
// If set to anything else, Git Town considers it the browser executable to use.
type BrowserExecutable string

func (self BrowserExecutable) Get() (executable string, useBrowser bool) { //nolint: nonamedreturns // the names really help understand the meaning of the return variables here
	if self == NoBrowser || self == "" {
		return "", false
	}
	return self.String(), true
}

func (self BrowserExecutable) String() string {
	return string(self)
}

func ParseBrowser(value, _ string) (Option[BrowserExecutable], error) {
	if value == "" {
		return None[BrowserExecutable](), nil
	}
	return Some(BrowserExecutable(value)), nil
}

func ParseBrowserHas(value string, has bool) (Option[BrowserExecutable], Option[BrowserEnabled], error) {
	if !has {
		return None[BrowserExecutable](), None[BrowserEnabled](), nil
	}
	if value == "" {
		return None[BrowserExecutable](), Some(BrowserEnabled(false)), nil
	}
	return Some(BrowserExecutable(value)), Some(BrowserEnabled(true)), nil
}
