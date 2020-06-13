package git

import "github.com/git-town/git-town/src/command"

// Root directory is cached in order to minimize the number of git commands run.
var rootDirectory string

// GetRootDirectory returns the path of the rood directory of the current repository,
// i.e. the directory that contains the ".git" folder.
func GetRootDirectory() string {
	if rootDirectory == "" {
		rootDirectory = command.MustRun("git", "rev-parse", "--show-toplevel").OutputSanitized()
	}
	return rootDirectory
}

// IsRebaseInProgress returns whether the local repository is in the middle of
// an unfinished rebase process.
func IsRebaseInProgress() bool {
	return command.MustRun("git", "status").OutputContainsText("rebase in progress")
}
