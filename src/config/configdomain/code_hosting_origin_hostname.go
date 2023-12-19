package configdomain

// GitLabToken is a bearer token to use with the GitLab API.
type CodeHostingOriginHostname string

func (self CodeHostingOriginHostname) String() string {
	return string(self)
}

func NewCodeHostingOriginHostnameRef(value string) *CodeHostingOriginHostname {
	token := CodeHostingOriginHostname(value)
	return &token
}
