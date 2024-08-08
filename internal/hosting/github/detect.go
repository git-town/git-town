package github

import "github.com/git-town/git-town/v15/internal/git/giturl"

// Detect indicates whether the current repository is hosted on a GitHub server.
func Detect(remoteURL giturl.Parts) bool {
	return remoteURL.Host == "github.com"
}
