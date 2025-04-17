package gitea

import "github.com/git-town/git-town/v19/internal/git/giturl"

// Detect indicates whether the current repository is hosted on a gitea server.
func Detect(remoteURL giturl.Parts) bool {
	return remoteURL.Host == "gitea.com"
}
