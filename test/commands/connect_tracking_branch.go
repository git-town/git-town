package commands

// ConnectTrackingBranch connects the branch with the given name to its counterpart at origin.
// The branch must exist.
func ConnectTrackingBranch(shell Shell, name string) error {
	_, err := shell.Run("git", "branch", "--set-upstream-to=origin/"+name, name)
	return err
}
