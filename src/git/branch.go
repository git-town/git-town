package git

import (
	"strings"

	"github.com/git-town/git-town/src/command"
	"github.com/git-town/git-town/src/util"
)

// GetLocalBranches returns the names of all branches in the local repository,
// ordered alphabetically.
func GetLocalBranches() (result []string) {
	for _, line := range command.MustRun("git", "branch").OutputLines() {
		line = strings.Trim(line, "* ")
		line = strings.TrimSpace(line)
		result = append(result, line)
	}
	return
}

// GetLocalBranchesWithoutMain returns the names of all branches in the local repository,
// ordered alphabetically without the main branch.
func GetLocalBranchesWithoutMain() (result []string) {
	mainBranch := Config().GetMainBranch()
	for _, branch := range GetLocalBranches() {
		if branch != mainBranch {
			result = append(result, branch)
		}
	}
	return
}

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

// HasLocalBranch returns whether the local repository contains
// a branch with the given name.
func HasLocalBranch(branchName string) bool {
	return util.DoesStringArrayContain(GetLocalBranches(), branchName)
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
