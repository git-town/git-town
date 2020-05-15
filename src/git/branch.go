package git

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/src/command"
	"github.com/git-town/git-town/src/util"
)

// DoesBranchHaveUnmergedCommits returns whether the branch with the given name
// contains commits that are not merged into the main branch
func DoesBranchHaveUnmergedCommits(branchName string) bool {
	return command.MustRun("git", "log", Config().GetMainBranch()+".."+branchName).OutputSanitized() != ""
}

// EnsureBranchInSync enforces that a branch with the given name is in sync with its tracking branch
func EnsureBranchInSync(branchName, errorMessageSuffix string) {
	util.Ensure(IsBranchInSync(branchName), fmt.Sprintf("%q is not in sync with its tracking branch. %s", branchName, errorMessageSuffix))
}

// EnsureDoesNotHaveBranch enforces that a branch with the given name does not exist
func EnsureDoesNotHaveBranch(branchName string) {
	util.Ensure(!HasBranch(branchName), fmt.Sprintf("A branch named %q already exists", branchName))
}

// EnsureHasBranch enforces that a branch with the given name exists
func EnsureHasBranch(branchName string) {
	util.Ensure(HasBranch(branchName), fmt.Sprintf("There is no branch named %q", branchName))
}

// EnsureHasLocalBranch enforces that a local branch with the given name exists
func EnsureHasLocalBranch(branchName string) {
	util.Ensure(HasLocalBranch(branchName), fmt.Sprintf("There is no local branch named %q", branchName))
}

// EnsureIsNotMainBranch enforces that a branch with the given name is not the main branch
func EnsureIsNotMainBranch(branchName, errorMessage string) {
	util.Ensure(!Config().IsMainBranch(branchName), errorMessage)
}

// EnsureIsNotPerennialBranch enforces that a branch with the given name is not a perennial branch
func EnsureIsNotPerennialBranch(branchName, errorMessage string) {
	util.Ensure(!Config().IsPerennialBranch(branchName), errorMessage)
}

// EnsureIsPerennialBranch enforces that a branch with the given name is a perennial branch
func EnsureIsPerennialBranch(branchName, errorMessage string) {
	util.Ensure(Config().IsPerennialBranch(branchName), errorMessage)
}

// GetExpectedPreviouslyCheckedOutBranch returns what is the expected previously checked out branch
// given the inputs
func GetExpectedPreviouslyCheckedOutBranch(initialPreviouslyCheckedOutBranch, initialBranch string) string {
	if HasLocalBranch(initialPreviouslyCheckedOutBranch) {
		if GetCurrentBranchName() == initialBranch || !HasLocalBranch(initialBranch) {
			return initialPreviouslyCheckedOutBranch
		}
		return initialBranch
	}
	return Config().GetMainBranch()
}

// GetLocalBranches returns the names of all branches in the local repository,
// ordered alphabetically
func GetLocalBranches() (result []string) {
	for _, line := range command.MustRun("git", "branch").OutputLines() {
		line = strings.Trim(line, "* ")
		line = strings.TrimSpace(line)
		result = append(result, line)
	}
	return
}

// GetLocalBranchesWithoutMain returns the names of all branches in the local repository,
// ordered alphabetically without the main branch
func GetLocalBranchesWithoutMain() (result []string) {
	mainBranch := Config().GetMainBranch()
	for _, branch := range GetLocalBranches() {
		if branch != mainBranch {
			result = append(result, branch)
		}
	}
	return
}

// GetLocalBranchesWithDeletedTrackingBranches returns the names of all branches
// whose remote tracking branches have been deleted
func GetLocalBranchesWithDeletedTrackingBranches() (result []string) {
	for _, line := range command.MustRun("git", "branch", "-vv").OutputLines() {
		line = strings.Trim(line, "* ")
		parts := strings.SplitN(line, " ", 2)
		branchName := parts[0]
		deleteTrackingBranchStatus := fmt.Sprintf("[%s: gone]", GetTrackingBranchName(branchName))
		if strings.Contains(parts[1], deleteTrackingBranchStatus) {
			result = append(result, branchName)
		}
	}
	return
}

// GetLocalBranchesWithMainBranchFirst returns the names of all branches
// that exist in the local repository,
// ordered to have the name of the main branch first,
// then the names of the branches, ordered alphabetically
func GetLocalBranchesWithMainBranchFirst() (result []string) {
	mainBranch := Config().GetMainBranch()
	result = append(result, mainBranch)
	for _, branch := range GetLocalBranches() {
		if branch != mainBranch {
			result = append(result, branch)
		}
	}
	return
}

// GetPreviouslyCheckedOutBranch returns the name of the previously checked out branch
func GetPreviouslyCheckedOutBranch() string {
	outcome, err := command.Run("git", "rev-parse", "--verify", "--abbrev-ref", "@{-1}")
	if err != nil {
		return ""
	}
	return outcome.OutputSanitized()
}

// GetTrackingBranchName returns the name of the remote branch
// that corresponds to the local branch with the given name
func GetTrackingBranchName(branchName string) string {
	return "origin/" + branchName
}

// HasBranch returns whether the repository contains a branch with the given name.
// The branch does not have to be present on the local repository.
func HasBranch(branchName string) bool {
	for _, line := range command.MustRun("git", "branch", "-a").OutputLines() {
		line = strings.Trim(line, "* ")
		line = strings.TrimSpace(line)
		line = strings.Replace(line, "remotes/origin/", "", 1)
		if line == branchName {
			return true
		}
	}
	return false
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

// IsBranchInSync returns whether the branch with the given name is in sync with its tracking branch
func IsBranchInSync(branchName string) bool {
	if HasTrackingBranch(branchName) {
		localSha := GetBranchSha(branchName)
		remoteSha := GetBranchSha(GetTrackingBranchName(branchName))
		return localSha == remoteSha
	}
	return true
}

// ShouldBranchBePushed returns whether the local branch with the given name
// contains commits that have not been pushed to the remote.
func ShouldBranchBePushed(branchName string) bool {
	trackingBranchName := GetTrackingBranchName(branchName)
	return command.MustRun("git", "rev-list", "--left-right", branchName+"..."+trackingBranchName).OutputSanitized() != ""
}

// Helpers

// Remote branches are cached in order to minimize the number of git commands run
var remoteBranches []string
var remoteBranchesInitialized bool

func getRemoteBranches() []string {
	if !remoteBranchesInitialized {
		remoteBranches = command.MustRun("git", "branch", "-r").OutputLines()
		remoteBranchesInitialized = true
	}
	return remoteBranches
}
