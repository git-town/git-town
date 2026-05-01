package browserdomain

import (
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

const NoBrowser = Browser("(none)")

// Browser indicates a custom browser to use.
// If set to "" or "(none)", browsers are disabled.
// If set to anything else, Git Town considers it the browser executable to use.
type Browser string

func (self Browser) Get() (executable string, useBrowser bool) { //nolint: nonamedreturns // the names really help understand the meaning of the return variables here
	if self == NoBrowser || self == "" {
		return "", false
	}
	return self.String(), true
}

func (self Browser) String() string {
	return string(self)
}

func NewBrowserErr(value, _ string) (Option[Browser], error) {
	return NewOption(Browser(value)), nil
}

func ParseBrowser(value string, has bool) Option[Browser] {
	if !has {
		return None[Browser]()
	}
	if value == "" {
		return Some(NoBrowser)
	}
	return Some(Browser(value))
}

func ParseBrowserErr(value string, has bool) (Option[Browser], error) {
	return NewBrowserErr(value, "")
}

// Indicates whether to use the browser.
func BrowserEnabled(setting Option[Browser]) bool {
	browser, hasBrowser := setting.Get()
	if !hasBrowser {
		return true
	}
	_, useBrowser := browser.Get()
	return useBrowser
}
