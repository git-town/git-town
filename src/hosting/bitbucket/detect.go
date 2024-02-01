package bitbucket

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/giturl"
)

// Detect indicates whether the current repository is hosted on a GitHub server.
func Detect(originURL *giturl.Parts, hostingPlatform configdomain.HostingPlatform) bool {
	return originURL != nil && (originURL.Host == "bitbucket.org" || hostingPlatform == configdomain.HostingPlatformBitbucket)
}
