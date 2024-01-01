package configdomain

type CodeHostingPlatformName string

func (self CodeHostingPlatformName) String() string {
	return string(self)
}

func NewCodeHostingPlatformNameRef(value string) *CodeHostingPlatformName {
	token := CodeHostingPlatformName(value)
	return &token
}
