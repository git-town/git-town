package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type HostingOriginHostname string

func (self HostingOriginHostname) String() string {
	return string(self)
}

func ParseHostingOriginHostname(value string) Option[HostingOriginHostname] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[HostingOriginHostname]()
	}
	return Some(HostingOriginHostname(value))
}
