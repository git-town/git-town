package configdomain

import (
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

type HostingOriginHostname stringss.Trimmed

func (self HostingOriginHostname) String() string {
	return string(self)
}

func ParseHostingOriginHostname(value stringss.Trimmed) Option[HostingOriginHostname] {
	if value == "" {
		return None[HostingOriginHostname]()
	}
	return Some(HostingOriginHostname(value))
}
