package gitea

import (
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/giturl"
	"github.com/git-town/git-town/v14/src/gohacks"
)

// Detect indicates whether the current repository is hosted on a GitHub server.
func Detect(originURL *giturl.Parts, hostingPlatform gohacks.Option[configdomain.HostingPlatform]) bool {
	if originURL != nil && originURL.Host == "gitea.com" {
		return true
	}
	if value, has := hostingPlatform.Get(); has {
		return value == configdomain.HostingPlatformGitea
	}
	return false
}
