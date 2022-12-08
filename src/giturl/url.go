// Package giturl provides facilities to work with the special URL formats used in Git remotes.
package giturl

import (
	"regexp"
)

// Parts contains recognized parts of a Git URL
type Parts struct {
	Protocol string // the protocol used to access the Git server, typically one of "ssh" or "https"
	User     string // optional username
	Host     string // hostname of the Git server
	Org      string // name of the organization that the repo is in
	Repo     string // name of the repository
}

func Parse(url string) *Parts {
	httpsRe := regexp.MustCompile(`https://(?P<user>.*@)?(?P<host>.*\/)(?P<org>.*\/)(?P<repo>.*)\.git`)
	matches := httpsRe.FindStringSubmatch(url)
	if matches != nil {
		return &Parts{
			Protocol: "https",
			User:     trimLast(matches[1]),
			Host:     trimLast(matches[2]),
			Org:      trimLast(matches[3]),
			Repo:     matches[4],
		}
	}
	sshRe := regexp.MustCompile(`(?P<user>.*@)?(?P<host>.*:)(?P<org>.*\/)(?P<repo>.*).git`)
	matches = sshRe.FindStringSubmatch(url)
	if matches != nil {
		return &Parts{
			Protocol: "ssh",
			User:     trimLast(matches[1]),
			Host:     trimLast(matches[2]),
			Org:      trimLast(matches[3]),
			Repo:     matches[4],
		}
	}
	return nil
}

// trims the last character of the given text, if the
func trimLast(text string) string {
	textLen := len(text)
	if textLen == 0 {
		return text
	}
	return text[:textLen-1]
}
