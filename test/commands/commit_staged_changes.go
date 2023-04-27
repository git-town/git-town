package commands

import "github.com/git-town/git-town/v8/test/subshell"

// CommitStagedChanges commits the currently staged changes.
func CommitStagedChanges(shell subshell.Mocking, message string) error {
	_, err := shell.Run("git", "commit", "-m", message)
	return err
}
