package git

import (
	"fmt"
	"os"

	"github.com/Originate/git-town/lib/util"
)

// EnsureDoesNotHaveConflicts asserts that the workspace
// has no unresolved merge conflicts.
func EnsureDoesNotHaveConflicts() {
	util.Ensure(!HasConflicts(), "You must resolve the conflicts before continuing")
}

// EnsureDoesNotHaveOpenChanges assets that the workspace
// has no open changes
func EnsureDoesNotHaveOpenChanges(message string) {
	util.Ensure(HasOpenChanges(), "You have uncommitted changes. "+message)
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

// HasShippableChanges returns whether the supplied branch has an changes
// not currently on the main branchName
func HasShippableChanges(branchName string) bool {
	return util.GetCommandOutput("git", "diff", GetMainBranch()+".."+branchName) != ""
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
