package dryrun

var currentBranchName = ""
var isActive = false

// Activate enables dry-run mode
func Activate(initialBranchName string) {
	isActive = true
	SetCurrentBranchName(initialBranchName)
}

// IsActive returns whether of not dry-run mode is active
func IsActive() bool {
	return isActive
}

// GetCurrentBranchName returns the current branch name for dry-run mode
func GetCurrentBranchName() string {
	return currentBranchName
}

// SetCurrentBranchName sets the current branch name for dry-run mode
func SetCurrentBranchName(branchName string) {
	currentBranchName = branchName
}
