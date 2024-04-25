package gitlab

import (
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/giturl"
)

// Detect indicates whether the current repository is hosted on a GitLab server.
func Detect(originURL *giturl.Parts, userOverride configdomain.HostingPlatform) bool {
	return originURL != nil && (originURL.Host == "gitlab.com" || userOverride == configdomain.HostingPlatformGitLab)
}
