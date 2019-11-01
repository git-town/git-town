package script

import "github.com/pkg/errors"

// SquashMerge squash merges the given branch into the current branch
func SquashMerge(branchName string) error {
	err := RunCommand("git", "merge", "--squash", branchName)
	if err != nil {
		return errors.Wrap(err, "cannot run git merge")
	}
	return nil
}
