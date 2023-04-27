package commands

func PushBranchToRemote(shell Shell, branch, remote string) error {
	_, err := shell.Run("git", "push", "-u", remote, branch)
	return err
}
