package git

import (
	"fmt"
	"os"

	"github.com/Originate/git-town/lib/util"
)

// EnsureDoesNotHaveConflicts asserts that the workspace
// has no unresolved merge conflicts.
func EnsureDoesNotHaveConflicts() {
	if HasConflicts() {
		util.ExitWithErrorMessage("You must resolve the conflicts before continuing")
	}
}

// GetRootDirectory returns the path of the rood directory of the current repository,
// i.e. the directory that contains the ".git" folder.
func GetRootDirectory() string {
	return util.GetCommandOutput("git", "rev-parse", "--show-toplevel")
}

// HasConflicts returns whether the local repository currently has unresolved merge conflicts.
func HasConflicts() bool {
	return util.DoesCommandOuputContain([]string{"git", "status"}, "Unmerged paths")
}

// HasOpenChanges returns whether the local repository contains uncommitted changes.
func HasOpenChanges() bool {
	return util.GetCommandOutput("git", "status", "--porcelain") != ""
}

// IsMergeInProgress returns whether the local repository is in the middle of
// an unfinished merge process.
func IsMergeInProgress() bool {
	_, err := os.Stat(fmt.Sprintf("%s/.git/MERGE_HEAD", GetRootDirectory()))
	return err == nil
}

// IsRebaseInProgress returns whether the local repository is in the middle of
// an unfinished rebase process.
func IsRebaseInProgress() bool {
	return util.DoesCommandOuputContain([]string{"git", "status"}, "rebase in progress")
}
