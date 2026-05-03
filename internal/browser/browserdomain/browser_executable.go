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

func ParseBrowserOpt(valueOpt Option[string], _ string) (Option[BrowserExecutable], error) {
	value, has := valueOpt.Get()
	if !has {
		return None[BrowserExecutable](), nil
	}
	if value == "" || value == NoBrowser.String() {
		return Some(NoBrowser), nil
	}
	return Some(BrowserExecutable(value)), nil
}

// Indicates whether to use the browser.
func BrowserEnabled(setting Option[BrowserExecutable]) bool {
	browser, hasBrowser := setting.Get()
	if !hasBrowser {
		return true
	}
	_, useBrowser := browser.Get()
	return useBrowser
}
