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

func NewCodeHostingPlatformRef(value string) *CodeHostingPlatform {
	token := CodeHostingPlatform(value)
	return &token
}
