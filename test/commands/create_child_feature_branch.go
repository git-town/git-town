package commands

// CreateChildFeatureBranch creates a branch with the given name and parent in this repository.
// The parent branch must already exist.
func CreateChildFeatureBranch(cmds *TestCommands, name string, parent string) error {
	err := CreateBranch(cmds, name, parent)
	if err != nil {
		return err
	}
	return cmds.Config.SetParent(name, parent)
}
