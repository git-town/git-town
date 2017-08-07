package drivers

// CodeHostingDriver defines the structure of drivers
// for the different code hosting services
type CodeHostingDriver struct {

	// CanBeUsed returns whether this driver can manage
	// a repository with the given hostname
	CanBeUsed func(hostname string) bool

	// GetNewPullRequestURL returns the URL of the page
	// to create a new pull request online
	GetNewPullRequestURL func(repository string, branch string, parentBranch string) string

	// GetRepositoryURL returns the URL where the given repository
	// can be found online
	GetRepositoryURL func(repository string) string

	// HostingServiceName contains the name of the code hosting service
	HostingServiceName string
}
