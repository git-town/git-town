package repo

// RemoveRemote deletes the Git remote with the given name.
func RemoveRemote(repo *Repo, name string) error {
	repo.Config().RemotesCache.Invalidate()
	_, err := repo.Run("git", "remote", "rm", name)
	return err
}
