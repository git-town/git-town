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
	switch value {
	case "", "auto-detect":
		return CodeHostingPlatformAutoDetect
	case "bitbucket", "BitBucket":
		return CodeHostingPlatformBitBucket
	case "gitea", "Gitea":
		return CodeHostingPlatformGitea
	case "github", "GitHub":
		return CodeHostingPlatformGitHub
	case "gitlab", "GitLab":
		return CodeHostingPlatformGitLab
	}
	panic("unknown code hosting platform: " + value)
}

func NewCodeHostingPlatformRef(value string) *CodeHostingPlatform {
	result := CodeHostingPlatform(value)
	return &result
}
