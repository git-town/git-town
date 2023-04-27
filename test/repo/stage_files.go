package repo

// StageFiles adds the file with the given name to the Git index.
func StageFiles(repo *Repo, names ...string) error {
	args := append([]string{"add"}, names...)
	_, err := repo.Run("git", args...)
	return err
}
