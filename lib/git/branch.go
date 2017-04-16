package git

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/Originate/git-town/lib/util"
)

// DoesBranchHaveUnmergedCommits returns whether the branch with the given name
// contains commits that are not merged into the main branch
func DoesBranchHaveUnmergedCommits(branchName string) bool {
	return util.GetCommandOutput("git", "log", GetMainBranch()+".."+branchName) != ""
}

// EnsureBranchInSync enforces that a branch with the given name is in sync with its tracking branch
func EnsureBranchInSync(branchName, errorMessageSuffix string) {
	util.Ensure(IsBranchInSync(branchName), fmt.Sprintf("'%s' is not in sync with its tracking branch. %s", branchName, errorMessageSuffix))
}

// EnsureDoesNotHaveBranch enforces that a branch with the given name does not exist
func EnsureDoesNotHaveBranch(branchName string) {
	util.Ensure(!HasBranch(branchName), fmt.Sprintf("A branch named '%s' already exists", branchName))
}

// EnsureHasBranch enforces that a branch with the given name exists
func EnsureHasBranch(branchName string) {
	util.Ensure(HasBranch(branchName), fmt.Sprintf("There is no branch named '%s'", branchName))
}

// EnsureIsNotMainBranch enforces that a branch with the given name is not the main branch
func EnsureIsNotMainBranch(branchName, errorMessage string) {
	util.Ensure(!IsMainBranch(branchName), errorMessage)
}

// EnsureIsNotPerennialBranch enforces that a branch with the given name is not a perennial branch
func EnsureIsNotPerennialBranch(branchName, errorMessage string) {
	util.Ensure(!IsPerennialBranch(branchName), errorMessage)
}

// EnsureIsPerennialBranch enforces that a branch with the given name is a perennial branch
func EnsureIsPerennialBranch(branchName, errorMessage string) {
	util.Ensure(IsPerennialBranch(branchName), errorMessage)
}

// GetCurrentBranchName returns the name of the currently checked out branch
func GetCurrentBranchName() string {
	if IsRebaseInProgress() {
		return getCurrentBranchNameDuringRebase()
	}
	return util.GetCommandOutput("git", "rev-parse", "--abbrev-ref", "HEAD")
}

// GetLocalBranches returns the names of all branches in the local repository,
// ordered alphabetically
func GetLocalBranches() (result []string) {
	output := util.GetCommandOutput("git", "branch")
	for _, line := range strings.Split(output, "\n") {
		line = strings.Trim(line, "* ")
		line = strings.TrimSpace(line)
		result = append(result, line)
	}
	return
}

// GetLocalBranchesWithDeletedTrackingBranches returns the names of all branches
// whose remote tracking branches have been deleted
func GetLocalBranchesWithDeletedTrackingBranches() (result []string) {
	output := util.GetCommandOutput("git", "branch", "-vv")
	for _, line := range strings.Split(output, "\n") {
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
	mainBranch := GetMainBranch()
	result = append(result, mainBranch)
	for _, branch := range GetLocalBranches() {
		if branch != mainBranch {
			result = append(result, branch)
		}
	}
	return
}

// GetTrackingBranchName returns the name of the remote branch
// that corresponds to the local branch with the given name
func GetTrackingBranchName(branchName string) string {
	return "origin/" + branchName
}

// HasBranch returns whether the repository contains a branch with the given name.
// The branch does not have to be present on the local repository.
func HasBranch(branchName string) bool {
	output := util.GetCommandOutput("git", "branch", "-a")
	for _, line := range strings.Split(output, "\n") {
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
	output := util.GetCommandOutput("git", "branch", "-r")
	for _, line := range strings.Split(output, "\n") {
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
	output := util.GetCommandOutput("git", "rev-list", "--left-right", branchName+"..."+trackingBranchName)
	return output != ""
}

// Helpers

func getCurrentBranchNameDuringRebase() string {
	filename := fmt.Sprintf("%s/.git/rebase-apply/head-name", GetRootDirectory())
	rawContent, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	content := strings.TrimSpace(string(rawContent))
	return strings.Replace(content, "refs/heads/", "", -1)
}
