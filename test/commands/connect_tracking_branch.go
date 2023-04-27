package commands

// ConnectTrackingBranch connects the branch with the given name to its counterpart at origin.
// The branch must exist.
func ConnectTrackingBranch(repo Repo, name string) error {
	_, err := repo.Run("git", "branch", "--set-upstream-to=origin/"+name, name)
	return err
}
