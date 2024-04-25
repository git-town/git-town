package gitea

import (
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/giturl"
)

// Detect indicates whether the current repository is hosted on a GitHub server.
func Detect(originURL *giturl.Parts, userOverride configdomain.HostingPlatform) bool {
	return originURL != nil && (originURL.Host == "gitea.com" || userOverride == configdomain.HostingPlatformGitea)
}
