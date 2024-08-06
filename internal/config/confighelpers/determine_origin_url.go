package confighelpers

import (
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/git/giturl"
	. "github.com/git-town/git-town/v15/internal/gohacks/prelude"
)

func DetermineOriginURL(originURL string, originOverride Option[configdomain.HostingOriginHostname]) Option[giturl.Parts] {
	url, hasURL := giturl.Parse(originURL).Get()
	if !hasURL {
		return None[giturl.Parts]()
	}
	if value, has := originOverride.Get(); has {
		url.Host = value.String()
	}
	return Some(url)
}
