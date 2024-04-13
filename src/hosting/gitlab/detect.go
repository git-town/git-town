package gitlab

import (
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/giturl"
)

// Detect indicates whether the current repository is hosted on a GitLab server.
func Detect(originURL *giturl.Parts, hostingPlatform configdomain.HostingPlatform) bool {
	return originURL != nil && (originURL.Host == "gitlab.com" || hostingPlatform == configdomain.HostingPlatformGitLab)
}
