package commands

func PushBranchToRemote(repo Repo, branch, remote string) error {
	_, err := repo.Run("git", "push", "-u", remote, branch)
	return err
}
