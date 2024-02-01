package hostingdomain

type Platform string

const (
	PlatformGitea     Platform = "Gitea"
	PlatformGithub    Platform = "GitHub"
	PlatformGitlab    Platform = "GitLab"
	PlatformBitbucket Platform = "Bitbucket"
	PlatformNone      Platform = "None"
)
