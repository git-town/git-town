package repo

// AddSubmodule adds a Git submodule with the given URL to this repository.
func AddSubmodule(repo *Repo, url string) error {
	_, err := repo.Run("git", "submodule", "add", url)
	if err != nil {
		return err
	}
	_, err = repo.Run("git", "commit", "-m", "added submodule")
	return err
}
