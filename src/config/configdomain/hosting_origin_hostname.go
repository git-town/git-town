package configdomain

type HostingOriginHostname string

func (self HostingOriginHostname) String() string {
	return string(self)
}

func NewHostingOriginHostnameRef(value string) *HostingOriginHostname {
	token := HostingOriginHostname(value)
	return &token
}
