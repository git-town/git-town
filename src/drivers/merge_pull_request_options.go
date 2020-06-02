package drivers

// MergePullRequestOptions defines the options to the MergePullRequest function.
type MergePullRequestOptions struct {
	Branch            string
	PullRequestNumber int
	CommitMessage     string
	LogRequests       bool
	ParentBranch      string
}
