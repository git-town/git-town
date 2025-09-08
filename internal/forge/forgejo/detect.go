package forgejo

import "github.com/git-town/git-town/v21/internal/git/giturl"

// Detect indicates whether the current repository is hosted on a Forgejo server.
func Detect(remoteURL giturl.Parts) bool {
	return remoteURL.Host == "codeberg.org"
}
