package commands

func PushBranch(repo *Repo) error {
	_, err := repo.Run("git", "push")
	return err
}
