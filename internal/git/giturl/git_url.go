// Package giturl provides facilities to work with the special URL formats used in Git remotes.
package giturl

import (
	"regexp"
	"strings"

	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Parts contains recognized parts of a Git URL.
type Parts struct {
	User       Option[string] // optional username
	Host       string         // hostname of the Git server
	Org        string         // name of the organization that the repo is in
	Repo       string         // name of the repository
	Supergroup Option[string] // optional over-arching grouping
}

func Parse(url string) Option[Parts] {
	pattern := `^` +
		// ignore transport protocol
		`(?:[^:]+://)?` +
		// capture "user@"
		`(?P<user>.*@)?` +
		// capture "host:" or "host/"
		`(?P<host>.*?)` +
		// ignore the port
		`(?:[:/]\[^/]*\/)?` +
		// capture "supergroup/org/repo" path
		`(.*)` +
		// ignore trailing ".git"
		`(?:\.git)?` +
		`$`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(url)
	if matches == nil {
		// NOTE: if we can't parse a Git URL, we simply ignore it.
		// This is because the URLs might be on the filesystem.
		// Remotes on the filesystem are not an error condition.
		return None[Parts]()
	}
	parts := strings.Split(matches[3], "/")
	var supergroup Option[string]
	var org string
	var repo string
	switch len(parts) {
	case 2:
		org = parts[0]
		repo = parts[1]
	case 3:
		supergroup = Some(parts[0])
		org = parts[1]
		repo = parts[2]
	}
	return Some(Parts{
		Host:       trimLast(matches[2]),
		Supergroup: supergroup,
		Org:        trimLast(org),
		Repo:       repo,
		User:       NewOption(trimLast(matches[1])),
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
