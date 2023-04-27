package commands

// Fetch retrieves the updates from the origin repo.
func Fetch(repo Repo) error {
	_, err := repo.Run("git", "fetch")
	return err
}
