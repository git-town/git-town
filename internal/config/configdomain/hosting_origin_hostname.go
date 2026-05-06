package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v23/pkg/prelude"
)

type HostingOriginHostname string

func (self HostingOriginHostname) String() string {
	return string(self)
}

func ParseHostingOriginHostname(valueOpt Option[string]) Option[HostingOriginHostname] {
	if value, has := valueOpt.Get(); has {
		return Some(HostingOriginHostname(strings.TrimSpace(value)))
	}
	return None[HostingOriginHostname]()
}
