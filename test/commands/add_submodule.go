package commands

import "github.com/git-town/git-town/v8/test/subshell"

// AddSubmodule adds a Git submodule with the given URL to this repository.
func AddSubmodule(shell subshell.Mocking, url string) error {
	_, err := shell.Run("git", "submodule", "add", url)
	if err != nil {
		return err
	}
	_, err = shell.Run("git", "commit", "-m", "added submodule")
	return err
}
