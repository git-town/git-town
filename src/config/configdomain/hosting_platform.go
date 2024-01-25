package configdomain

type HostingPlatform string

const (
	HostingPlatformAutoDetect = "auto-detect"
	HostingPlatformBitBucket  = "bitbucket"
	HostingPlatformGitea      = "gitea"
	HostingPlatformGitHub     = "github"
	HostingPlatformGitLab     = "gitlab"
)

func (self HostingPlatform) String() string {
	return string(self)
}

func NewHostingPlatformRef(value string) *HostingPlatform {
	token := HostingPlatform(value)
	return &token
}
