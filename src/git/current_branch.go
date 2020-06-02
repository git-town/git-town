package git

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/git-town/git-town/src/command"
	"github.com/git-town/git-town/src/dryrun"
)

// The current branch in cached in order to minimize the number of git commands run.
var currentBranchCache string

// GetCurrentBranchName returns the name of the currently checked out branch.
func GetCurrentBranchName() string {
	if dryrun.IsActive() {
		return dryrun.GetCurrentBranchName()
	}
	if currentBranchCache == "" {
		if IsRebaseInProgress() {
			currentBranchCache = getCurrentBranchNameDuringRebase()
		} else {
			currentBranchCache = command.MustRun("git", "rev-parse", "--abbrev-ref", "HEAD").OutputSanitized()
		}
	}
	return currentBranchCache
}

// ClearCurrentBranchCache clears the cache of the current branch.
// This should be called when a rebase fails.
func ClearCurrentBranchCache() {
	currentBranchCache = ""
}

// UpdateCurrentBranchCache clears the cache of the current branch.
// This should be called when a new branch is checked out.
func UpdateCurrentBranchCache(branchName string) {
	currentBranchCache = branchName
}

// Helpers

func getCurrentBranchNameDuringRebase() string {
	rawContent, err := ioutil.ReadFile(fmt.Sprintf("%s/.git/rebase-apply/head-name", GetRootDirectory()))
	if err != nil {
		// Git 2.26 introduces a new rebase backend, see https://github.com/git/git/blob/master/Documentation/RelNotes/2.26.0.txt
		rawContent, err = ioutil.ReadFile(fmt.Sprintf("%s/.git/rebase-merge/head-name", GetRootDirectory()))
		if err != nil {
			panic(err)
		}
	}
	content := strings.TrimSpace(string(rawContent))
	return strings.Replace(content, "refs/heads/", "", -1)
}
