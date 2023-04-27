package commands

// CommitStagedChanges commits the currently staged changes.
func CommitStagedChanges(repo *Repo, message string) error {
	_, err := repo.Run("git", "commit", "-m", message)
	return err
}
