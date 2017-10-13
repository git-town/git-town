package drivers

// CodeHostingDriver defines the structure of drivers
// for the different code hosting services
type CodeHostingDriver interface {

	// CanBeUsed returns whether this driver can manage
	// a repository with the given hostname
	CanBeUsed(driverType string) bool

	// CanMergePullRequest returns whether or not MergePullRequest should be
	// called when shipping. If true, also returns the default commit message
	CanMergePullRequest(branch, parentBranch string) (bool, string, error)

	// GetNewPullRequestURL returns the URL of the page
	// to create a new pull request online
	GetNewPullRequestURL(branch, parentBranch string) string

	// MergePullRequest merges the pull request through the hosting service api
	MergePullRequest(MergePullRequestOptions) (string, error)

	// GetRepositoryURL returns the URL where the given repository
	// can be found online
	GetRepositoryURL() string

	// HostingServiceName returns the name of the code hosting service
	HostingServiceName() string

	// SetOriginURL configures the driver with the origin URL of the Git repo
	SetOriginURL(originURL string)

	// SetOriginHostname configures the driver with the origin hostname of the Git repo
	SetOriginHostname(originHostname string)

	// GetAPITokenKey returns the git config key value that the API token is stored under
	GetAPITokenKey() string

	// SetAPIToken configures the driver with API token
	SetAPIToken(apiToken string)
}
