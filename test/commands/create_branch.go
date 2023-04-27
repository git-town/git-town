package commands

import "github.com/git-town/git-town/v8/test/subshell"

// CreateBranch creates a new branch with the given name.
// The created branch is a normal branch.
// To create feature branches, use CreateFeatureBranch.
func CreateBranch(shell subshell.Mocking, name, parent string) error {
	_, err := shell.Run("git", "branch", name, parent)
	return err
}
