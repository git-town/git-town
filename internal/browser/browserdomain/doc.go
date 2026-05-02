// Package browserdomain provides the domain model for browsers.
package browserdomain

import . "github.com/git-town/git-town/v22/pkg/prelude"

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
