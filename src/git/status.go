package git

import (
	"github.com/git-town/git-town/src/command"
	"github.com/git-town/git-town/src/util"
)

// EnsureDoesNotHaveConflicts asserts that the workspace
// has no unresolved merge conflicts.
func EnsureDoesNotHaveConflicts() {
	util.Ensure(!HasConflicts(), "You must resolve the conflicts before continuing")
}

// EnsureDoesNotHaveOpenChanges assets that the workspace
// has no open changes
func EnsureDoesNotHaveOpenChanges(message string) {
	util.Ensure(!HasOpenChanges(), "You have uncommitted changes. "+message)
}

// Root directory is cached in order to minimize the number of git commands run
var rootDirectory string

// GetRootDirectory returns the path of the rood directory of the current repository,
// i.e. the directory that contains the ".git" folder.
func GetRootDirectory() string {
	if rootDirectory == "" {
		rootDirectory = command.MustRun("git", "rev-parse", "--show-toplevel").OutputSanitized()
	}
	return rootDirectory
}

// HasConflicts returns whether the local repository currently has unresolved merge conflicts.
func HasConflicts() bool {
	return command.MustRun("git", "status").OutputContainsText("Unmerged paths")
}

// HasOpenChanges returns whether the local repository contains uncommitted changes.
func HasOpenChanges() bool {
	return command.MustRun("git", "status", "--porcelain").OutputSanitized() != ""
}

// IsRebaseInProgress returns whether the local repository is in the middle of
// an unfinished rebase process.
func IsRebaseInProgress() bool {
	return command.MustRun("git", "status").OutputContainsText("rebase in progress")
}
