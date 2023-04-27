package commands

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func CreateBranch(repo Repo, name, parent string) error {
	_, err := repo.Run("git", "branch", name, parent)
	return err
}
