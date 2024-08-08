// Package giturl provides facilities to work with the special URL formats used in Git remotes.
package giturl

import (
	"regexp"

	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
)

// Parts contains recognized parts of a Git URL.
type Parts struct {
	User Option[string] // optional username
	Host string         // hostname of the Git server
	Org  string         // name of the organization that the repo is in
	Repo string         // name of the repository
}

func Parse(url string) Option[Parts] {
	pattern := `^` +
		// ignore transport protocol
		`(?:[^:]+://)?` +
		// capture "user@"
		`(?P<user>.*@)?` +
		// capture "host:" or "host/"
		`(?P<host>.*?[:/])` +
		// ignore the port
		`(?:\d+\/)?` +
		// capture "org/"
		`(?P<org>.*\/)` +
		// capture "repo"
		`(?P<repo>.*?)` +
		// ignore trailing ".git"
		`(?:\.git)?$`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(url)
	if matches == nil {
		// NOTE: if we can't parse a Git URL, we simply ignore it.
		// This is because the URLs might be on the filesystem.
		// Remotes on the filesystem are not an error condition.
		return None[Parts]()
	}
	return Some(Parts{
		Host: trimLast(matches[2]),
		Org:  trimLast(matches[3]),
		Repo: matches[4],
		User: NewOption(trimLast(matches[1])),
	})
}

// trimLast trims the last character of the given text.
// Handles empty strings gracefully.
func trimLast(text string) string {
	textLen := len(text)
	if textLen == 0 {
		return text
	}
	return text[:textLen-1]
}
