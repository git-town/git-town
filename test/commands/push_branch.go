package commands

func PushBranch(shell Shell) error {
	_, err := shell.Run("git", "push")
	return err
}
