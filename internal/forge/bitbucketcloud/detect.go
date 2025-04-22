package bitbucketcloud

import "github.com/git-town/git-town/v19/internal/git/giturl"

// Detect indicates whether the current repository is hosted on a Bitbucket server.
func Detect(remoteURL giturl.Parts) bool {
	return remoteURL.Host == "bitbucket.org"
}
