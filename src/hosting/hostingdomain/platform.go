package hostingdomain

type Platform string

const (
	PlatformGitea     Platform = "Gitea"
	PlatformGithub    Platform = "GitHub"
	PlatformGitlab    Platform = "GitLab"
	PlatformBitbucket Platform = "Bitbucket"
)

func (self Platform) String() string {
	return string(self)
}
