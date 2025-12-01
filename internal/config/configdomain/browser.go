package configdomain

import (
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type Browser string

func (self Browser) String() string {
	return string(self)
}

func (self Browser) NoBrowser() bool {
	return self == NoBrowser
}

func ParseBrowser(value, _ string) (Option[Browser], error) {
	if value == "" {
		return None[Browser](), nil
	}
	return Some(Browser(value)), nil
}

const NoBrowser = "(none)"
