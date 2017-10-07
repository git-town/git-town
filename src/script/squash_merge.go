package script

import "github.com/Originate/git-town/src/exit"

// SquashMerge squash merges the given branch into the current branch
func SquashMerge(branchName string) {
	err := RunCommand("git", "merge", "--squash", branchName)
	exit.OnWrap(err, "Error squash merging")
}
