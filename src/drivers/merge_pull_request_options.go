package drivers

// MergePullRequestOptions defines the options to the MergePullRequest function
type MergePullRequestOptions struct {
	Branch        string
	CommitMessage string
	LogRequests   bool
	ParentBranch  string
}
