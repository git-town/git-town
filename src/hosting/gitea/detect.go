package gitea

import (
	"github.com/git-town/git-town/v13/src/config/configdomain"
	"github.com/git-town/git-town/v13/src/git/giturl"
)

// Detect indicates whether the current repository is hosted on a GitHub server.
func Detect(originURL *giturl.Parts, hostingPlatform configdomain.HostingPlatform) bool {
	return originURL != nil && (originURL.Host == "gitea.com" || hostingPlatform == configdomain.HostingPlatformGitea)
}
