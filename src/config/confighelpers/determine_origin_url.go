package confighelpers

import (
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/giturl"
	"github.com/git-town/git-town/v14/src/gohacks"
)

func DetermineOriginURL(originURL string, originOverride gohacks.Option[configdomain.HostingOriginHostname], originURLCache configdomain.OriginURLCache) *giturl.Parts {
	cached, has := originURLCache[originURL]
	if has {
		return cached
	}
	url := giturl.Parse(originURL)
	if value, has := originOverride.Get(); has {
		url.Host = value.String()
	}
	originURLCache[originURL] = url
	return url
}
