package commands

// Fetch retrieves the updates from the origin repo.
func Fetch(shell Shell) error {
	_, err := shell.Run("git", "fetch")
	return err
}
