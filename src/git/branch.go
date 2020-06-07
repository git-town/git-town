package git

import (
	"strings"

	"github.com/git-town/git-town/src/command"
)

// GetPreviouslyCheckedOutBranch returns the name of the previously checked out branch.
func GetPreviouslyCheckedOutBranch() string {
	outcome, err := command.Run("git", "rev-parse", "--verify", "--abbrev-ref", "@{-1}")
	if err != nil {
		return ""
	}
	return outcome.OutputSanitized()
}

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

// IsBranchInSync returns whether the branch with the given name is in sync with its tracking branch.
func IsBranchInSync(branchName string) bool {
	if HasTrackingBranch(branchName) {
		localSha := GetBranchSha(branchName)
		remoteSha := GetBranchSha(GetTrackingBranchName(branchName))
		return localSha == remoteSha
	}
	return true
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
