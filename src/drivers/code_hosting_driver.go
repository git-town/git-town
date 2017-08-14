package drivers

// CodeHostingDriver defines the structure of drivers
// for the different code hosting services
type CodeHostingDriver interface {

	// CanBeUsed returns whether this driver can manage
	// a repository with the given hostname
	CanBeUsed() bool

	// GetNewPullRequestURL returns the URL of the page
	// to create a new pull request online
	GetNewPullRequestURL(repository string, branch string, parentBranch string) string

	// GetRepositoryURL returns the URL where the given repository
	// can be found online
	GetRepositoryURL(repository string) string

	// HostingServiceName returns the name of the code hosting service
	HostingServiceName() string

	// SetOriginURL configures the driver with the origin URL of the Git repo
	SetOriginURL(originURL string)
}
