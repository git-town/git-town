package commands

// CreateChildFeatureBranch creates a branch with the given name and parent in this repository.
// The parent branch must already exist.
func CreateChildFeatureBranch(repo Repo, name string, parent string) error {
	err := CreateBranch(repo, name, parent)
	if err != nil {
		return err
	}
	return repo.Config().SetParent(name, parent)
}
