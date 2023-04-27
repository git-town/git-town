package commands

// RemoveBranch deletes the branch with the given name from this repo.
func RemoveBranch(repo *Repo, name string) error {
	_, err := repo.Run("git", "branch", "-D", name)
	return err
}
