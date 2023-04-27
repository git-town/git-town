package commands

// CreateTag creates a tag with the given name.
func CreateTag(shell Shell, name string) error {
	_, err := shell.Run("git", "tag", "-a", name, "-m", name)
	return err
}
