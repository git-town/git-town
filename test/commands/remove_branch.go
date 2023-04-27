package commands

// RemoveBranch deletes the branch with the given name from this repo.
func RemoveBranch(shell Shell, name string) error {
	_, err := shell.Run("git", "branch", "-D", name)
	return err
}
