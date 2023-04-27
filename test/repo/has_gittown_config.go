package repo

// HasGitTownConfigNow indicates whether this repository contain Git Town specific configuration.
func HasGitTownConfigNow(repo *Repo) bool {
	output, err := repo.Run("git", "config", "--local", "--get-regex", "git-town")
	if err != nil {
		return false
	}
	return output != ""
}
