package configdomain

type CodeHostingPlatform string

const (
	CodeHostingPlatformAutoDetect CodeHostingPlatform = "auto-detect"
	CodeHostingPlatformBitBucket  CodeHostingPlatform = "bitbucket"
	CodeHostingPlatformGitea      CodeHostingPlatform = "gitea"
	CodeHostingPlatformGitHub     CodeHostingPlatform = "github"
	CodeHostingPlatformGitLab     CodeHostingPlatform = "gitlab"
)

func (self CodeHostingPlatform) String() string {
	return string(self)
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
	token := CodeHostingPlatform(value)
	return &token
}
