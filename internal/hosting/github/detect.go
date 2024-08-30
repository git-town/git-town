package github

import "github.com/git-town/git-town/v16/internal/git/giturl"

// Detect indicates whether the current repository is hosted on a GitHub server.
func Detect(remoteURL giturl.Parts) bool {
	return remoteURL.Host == "github.com"
}
