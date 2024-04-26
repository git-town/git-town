package confighelpers

import (
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/giturl"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

func DetermineOriginURL(originURL string, originOverride Option[configdomain.HostingOriginHostname], originURLCache configdomain.OriginURLCache) (giturl.Parts, error) {
	cached, has := originURLCache[originURL]
	if has {
		return cached, nil
	}

	url, err := giturl.Parse(originURL)
	if err != nil {
		var p giturl.Parts
		return p, err
	}
	if value, has := originOverride.Get(); has {
		url.Host = value.String()
	}
	originURLCache[originURL] = url
	return url, nil
}
