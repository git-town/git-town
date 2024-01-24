package configdomain

type CodeHostingPlatformName string

const (
	CodeHostingPlatformNameAutoDetect = "auto-detect"
	CodeHostingPlatformBitBucket      = "bitbucket"
	CodeHostingPlatformGitea          = "gitea"
	CodeHostingPlatformGitHub         = "github"
	CodeHostingPlatformGitLab         = "gitlab"
)

func (self CodeHostingPlatformName) String() string {
	return string(self)
}

func AllCodeHostingPlatformNames() []CodeHostingPlatformName {
	return []CodeHostingPlatformName{
		CodeHostingPlatformNameAutoDetect,
		CodeHostingPlatformBitBucket,
		CodeHostingPlatformGitea,
		CodeHostingPlatformGitHub,
		CodeHostingPlatformGitLab,
	}
}

func NewCodeHostingPlatformName(value string) CodeHostingPlatformName {
	for _, codeHostingPlatformName := range AllCodeHostingPlatformNames() {
		if codeHostingPlatformName.String() == value {
			return codeHostingPlatformName
		}
	}
	panic("unknown code hosting platform name: " + value)
}

func NewCodeHostingPlatformNameRef(value string) *CodeHostingPlatformName {
	result := NewCodeHostingPlatformName(value)
	return &result
}
