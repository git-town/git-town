// Package browserdomain provides the domain model for browsers.
package browserdomain

import . "github.com/git-town/git-town/v23/pkg/prelude"

func ParseBrowserOpt(valueOpt Option[string]) (Option[BrowserExecutable], Option[BrowserEnabled], error) {
	value, has := valueOpt.Get()
	if !has {
		return None[BrowserExecutable](), None[BrowserEnabled](), nil
	}
	if value == "" || value == NoBrowser.String() {
		return None[BrowserExecutable](), Some(BrowserEnabled(false)), nil
	}
	return Some(BrowserExecutable(value)), None[BrowserEnabled](), nil
}
