// Package giturl provides facilities to work with the special URL formats used in Git remotes.
package giturl

import (
	"regexp"
)

// Parts contains recognized parts of a Git URL.
type Parts struct {
	User string // optional username
	Host string // hostname of the Git server
	Org  string // name of the organization that the repo is in
	Repo string // name of the repository
}

func Parse(url string) *Parts {
	pattern := `^` +
		// ignore transport protocol
		`(?:.*?://)?` +
		// capture "user@"
		`(?P<user>.*@)?` +
		// capture "host:" or "hostname/"
		`(?P<host>.*?[:/])` +
		// capture "org/"
		`(?P<org>.*\/)` +
		// capture "repo$" or "repo.git$" and ignore the trailing ".git"
		`(?P<repo>.*?)?(?:\.git)?$`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(url)
	if matches != nil {
		return &Parts{
			User: trimLast(matches[1]),
			Host: trimLast(matches[2]),
			Org:  trimLast(matches[3]),
			Repo: matches[4],
		}
	}
	return nil
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
