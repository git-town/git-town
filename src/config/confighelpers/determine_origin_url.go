package confighelpers

import (
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/giturl"
)

func DetermineOriginURL(originURL string, originOverride configdomain.HostingOriginHostname, originURLCache configdomain.OriginURLCache) *giturl.Parts {
	cached, has := originURLCache[originURL]
	if has {
		return cached
	}
	url := giturl.Parse(originURL)
	if originOverride != "" {
		url.Host = string(originOverride)
	}
	originURLCache[originURL] = url
	return url
}
