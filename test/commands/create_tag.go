package commands

import "github.com/git-town/git-town/v8/test/subshell"

// CreateTag creates a tag with the given name.
func CreateTag(shell subshell.Mocking, name string) error {
	_, err := shell.Run("git", "tag", "-a", name, "-m", name)
	return err
}
