package script

import "fmt"

// SquashMerge squash merges the given branch into the current branch
func SquashMerge(branchName string) error {
	err := RunCommand("git", "merge", "--squash", branchName)
	if err != nil {
		return fmt.Errorf("cannot run git merge: %w", err)
	}
	return nil
}
