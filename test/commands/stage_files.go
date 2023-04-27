package commands

// StageFiles adds the file with the given name to the Git index.
func StageFiles(shell Shell, names ...string) error {
	args := append([]string{"add"}, names...)
	_, err := shell.Run("git", args...)
	return err
}
