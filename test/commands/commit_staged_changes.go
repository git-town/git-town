package commands

// CommitStagedChanges commits the currently staged changes.
func CommitStagedChanges(shell Shell, message string) error {
	_, err := shell.Run("git", "commit", "-m", message)
	return err
}
