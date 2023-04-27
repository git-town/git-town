package commands

// CreateTag creates a tag with the given name.
func CreateTag(repo Repo, name string) error {
	_, err := repo.Run("git", "tag", "-a", name, "-m", name)
	return err
}
