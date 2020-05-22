package drivers

// CodeHostingDriver defines the structure of drivers
// for the different code hosting services
type CodeHostingDriver interface {

	// WasActivated returns whether this driver was applicable
	// and activated for any given DriverOptions
	WasActivated(opts DriverOptions) bool

	// CanMergePullRequest returns whether or not MergePullRequest should be
	// called when shipping. If true, also returns the default commit message
	CanMergePullRequest(branch, parentBranch string) (canMerge bool, defaultCommitMessage string, pullRequestNumber int64, err error)

	// GetNewPullRequestURL returns the URL of the page
	// to create a new pull request online
	GetNewPullRequestURL(branch, parentBranch string) string

	// MergePullRequest merges the pull request through the hosting service api
	MergePullRequest(MergePullRequestOptions) (mergeSha string, err error)

	// GetRepositoryURL returns the URL where the given repository
	// can be found online
	GetRepositoryURL() string

	// HostingServiceName returns the name of the code hosting service
	HostingServiceName() string
}
