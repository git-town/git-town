package confighelpers

import (
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/git/giturl"
	. "github.com/git-town/git-town/v17/pkg/prelude"
)

func DetermineRemoteURL(urlStr string, override Option[configdomain.HostingOriginHostname]) Option[giturl.Parts] {
	url, hasURL := giturl.Parse(urlStr).Get()
	if !hasURL {
		return None[giturl.Parts]()
	}
	if value, has := override.Get(); has {
		url.Host = value.String()
	}
	return Some(url)
}
