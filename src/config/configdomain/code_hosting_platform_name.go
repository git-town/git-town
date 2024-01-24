package configdomain

type CodeHostingPlatform string

const (
	CodeHostingPlatformAutoDetect = "auto-detect"
	CodeHostingPlatformBitBucket  = "bitbucket"
	CodeHostingPlatformGitea      = "gitea"
	CodeHostingPlatformGitHub     = "github"
	CodeHostingPlatformGitLab     = "gitlab"
)

func (self CodeHostingPlatform) String() string {
	return string(self)
}

func AllCodeHostingPlatforms() []CodeHostingPlatform {
	return []CodeHostingPlatform{
		CodeHostingPlatformAutoDetect,
		CodeHostingPlatformBitBucket,
		CodeHostingPlatformGitea,
		CodeHostingPlatformGitHub,
		CodeHostingPlatformGitLab,
	}
}

func NewCodeHostingPlatform(value string) CodeHostingPlatform {
	for _, codeHostingPlatformName := range AllCodeHostingPlatforms() {
		if codeHostingPlatformName.String() == value {
			return codeHostingPlatformName
		}
	}
	panic("unknown code hosting platform: " + value)
}

func NewCodeHostingPlatformRef(value string) *CodeHostingPlatform {
	result := NewCodeHostingPlatform(value)
	return &result
}
