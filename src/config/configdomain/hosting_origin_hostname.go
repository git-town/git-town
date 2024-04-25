package configdomain

import "github.com/git-town/git-town/v14/src/gohacks"

type HostingOriginHostname string

func (self HostingOriginHostname) String() string {
	return string(self)
}

func NewHostingOriginHostname(value string) HostingOriginHostname {
	return HostingOriginHostname(value)
}

func NewHostingOriginHostnameOption(value string) gohacks.Option[HostingOriginHostname] {
	if value == "" {
		return gohacks.NewOptionNone[HostingOriginHostname]()
	}
	return gohacks.NewOption(NewHostingOriginHostname(value))
}
