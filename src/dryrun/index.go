package dryrun

import "github.com/git-town/git-town/src/cli"

var currentBranchName = ""
var isActive = false

// Activate causes all commands to not be run.
func Activate(currentBranch string) {
	cli.PrintDryRunMessage()
	isActive = true
	SetCurrentBranchName(currentBranch)
}

// IsActive returns whether of not dry-run mode is active.
func IsActive() bool {
	return isActive
}

// GetCurrentBranchName returns the current branch name for dry-run mode.
func GetCurrentBranchName() string {
	return currentBranchName
}

// SetCurrentBranchName sets the current branch name for dry-run mode.
func SetCurrentBranchName(branchName string) {
	currentBranchName = branchName
}
