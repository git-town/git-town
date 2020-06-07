package dryrun

import (
	"github.com/fatih/color"
)

var currentBranchName = ""
var isActive = false

// Activate enables dry-run mode.
func Activate(initialBranchName string) {
	isActive = true
	SetCurrentBranchName(initialBranchName)
}

// ActivateDryRun causes all commands to not be run.
func ActivateDryRun(currentBranch string) {
	_, err := color.New(color.FgBlue).Print(dryRunMessage)
	if err != nil {
		panic(err)
	}
	Activate(currentBranch)
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

var dryRunMessage = `
In dry run mode. No commands will be run. When run in normal mode, the command
output will appear beneath the command. Some commands will only be run if
necessary. For example: 'git push' will run if and only if there are local
commits not on the remote.
`
