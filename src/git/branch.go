package git

// GetTrackingBranchName returns the name of the remote branch
// that corresponds to the local branch with the given name.
func GetTrackingBranchName(branchName string) string {
	return "origin/" + branchName
}
