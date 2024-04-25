package bitbucket

import (
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/giturl"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

// Detect indicates whether the current repository is hosted on a GitHub server.
func Detect(originURL *giturl.Parts, hostingPlatform Option[configdomain.HostingPlatform]) bool {
	if originURL != nil && originURL.Host == "bitbucket.org" {
		return true
	}
	if value, has := hostingPlatform.Get(); has {
		return value == configdomain.HostingPlatformBitbucket
	}
	return false
}
