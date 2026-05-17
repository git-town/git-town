package browserdomain

import "github.com/git-town/git-town/v23/internal/gohacks/stringss"

const NoBrowser = BrowserExecutable("(none)")

// BrowserExecutable indicates a custom browser to use.
// If set to "" or "(none)", browsers are disabled.
// If set to anything else, Git Town considers it the browser executable to use.
type BrowserExecutable stringss.TrimmedString

func (self BrowserExecutable) Get() (executable string, useBrowser bool) { //nolint: nonamedreturns // the names really help understand the meaning of the return variables here
	if self == NoBrowser || self == "" {
		return "", false
	}
	return self.String(), true
}

func (self BrowserExecutable) String() string {
	return string(self)
}
