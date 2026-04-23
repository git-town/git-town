package browserdomain

import (
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

const NoBrowser = Browser("(none)")

// Browser indicates a custom browser to use.
// If set to "" or "(none)", browsers are disabled.
// IF set to anything else, Git Town considers it the browser executable to use.
type Browser string

func (self Browser) Get() (executable string, useBrowser bool) {
	if self == NoBrowser {
		return "", false
	}
	return self.String(), true
}

func (self Browser) String() string {
	return string(self)
}

func ParseBrowser(value, _ string) (Option[Browser], error) {
	if value == "" {
		return None[Browser](), nil
	}
	return Some(Browser(value)), nil
}

func ParseBrowserHas(value string, has bool) (Option[Browser], error) {
	if !has {
		return None[Browser](), nil
	}
	if value == "" {
		return Some(NoBrowser), nil
	}
	return Some(Browser(value)), nil
}

// Indicates whether to use the browser, and if so, the executable to use.
func UseBrowser(setting Option[Browser]) bool {
	browser, hasBrowser := setting.Get()
	if !hasBrowser {
		return true
	}
	_, useBrowser := browser.Get()
	return useBrowser
}
