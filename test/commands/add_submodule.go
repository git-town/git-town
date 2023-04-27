package commands

// AddSubmodule adds a Git submodule with the given URL to this repository.
func AddSubmodule(shell Shell, url string) error {
	_, err := shell.Run("git", "submodule", "add", url)
	if err != nil {
		return err
	}
	_, err = shell.Run("git", "commit", "-m", "added submodule")
	return err
}
