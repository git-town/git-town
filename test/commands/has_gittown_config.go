package commands

// HasGitTownConfigNow indicates whether this repository contain Git Town specific configuration.
func HasGitTownConfigNow(shell Shell) bool {
	output, err := shell.Run("git", "config", "--local", "--get-regex", "git-town")
	if err != nil {
		return false
	}
	return output != ""
}
