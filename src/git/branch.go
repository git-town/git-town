package git

import (
	"strings"

	"github.com/git-town/git-town/src/command"
)

// GetTrackingBranchName returns the name of the remote branch
// that corresponds to the local branch with the given name.
func GetTrackingBranchName(branchName string) string {
	return "origin/" + branchName
}

// HasTrackingBranch returns whether the local branch with the given name
// has a tracking branch.
func HasTrackingBranch(branchName string) bool {
	trackingBranchName := GetTrackingBranchName(branchName)
	for _, line := range getRemoteBranches() {
		if strings.TrimSpace(line) == trackingBranchName {
			return true
		}
	}
	return false
}

// Helpers

// Remote branches are cached in order to minimize the number of git commands run.
var remoteBranches []string
var remoteBranchesInitialized bool

func getRemoteBranches() []string {
	if !remoteBranchesInitialized {
		remoteBranches = command.MustRun("git", "branch", "-r").OutputLines()
		remoteBranchesInitialized = true
	}
	return remoteBranches
}
