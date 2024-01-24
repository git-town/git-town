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

func NewCodeHostingPlatformNameRef(value string) *CodeHostingPlatformName {
	token := CodeHostingPlatformName(value)
	return &token
}
