package configdomain

import (
	"strings"

	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
)

type HostingOriginHostname string

func (self HostingOriginHostname) String() string {
	return string(self)
}

func NewHostingOriginHostname(value string) HostingOriginHostname {
	value = strings.TrimSpace(value)
	if value == "" {
		panic("empty hosting origin hostname")
	}
	return HostingOriginHostname(value)
}

func NewHostingOriginHostnameOption(value string) Option[HostingOriginHostname] {
	value = strings.TrimSpace(value)
	if value == "" {
		return None[HostingOriginHostname]()
	}
	return Some(NewHostingOriginHostname(value))
}
