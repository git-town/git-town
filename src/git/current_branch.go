package git

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/command"
	"github.com/Originate/git-town/src/dryrun"
)

// The current branch in cached in order to minimize the number of git commands run
var currentBranchCache string

// GetCurrentBranchName returns the name of the currently checked out branch
func GetCurrentBranchName() string {
	if dryrun.IsActive() {
		return dryrun.GetCurrentBranchName()
	}
	if currentBranchCache == "" {
		if IsRebaseInProgress() {
			currentBranchCache = getCurrentBranchNameDuringRebase()
		} else {
			currentBranchCache = command.New("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
		}
	}
	return currentBranchCache
}

// ClearCurrentBranchCache clears the cache of the current branch.
// This should be called when a rebase fails
func ClearCurrentBranchCache() {
	currentBranchCache = ""
}

// UpdateCurrentBranchCache clears the cache of the current branch.
// This should be called when a new branch is checked out
func UpdateCurrentBranchCache(branchName string) {
	currentBranchCache = branchName
}

// Helpers

func getCurrentBranchNameDuringRebase() string {
	filename := fmt.Sprintf("%s/.git/rebase-apply/head-name", GetRootDirectory())
	rawContent, err := ioutil.ReadFile(filename)
	exit.If(err)
	content := strings.TrimSpace(string(rawContent))
	return strings.Replace(content, "refs/heads/", "", -1)
}
