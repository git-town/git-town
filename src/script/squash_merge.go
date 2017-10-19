package script

import "github.com/Originate/exit"

// SquashMerge squash merges the given branch into the current branch
func SquashMerge(branchName string) {
	err := RunCommand("git", "merge", "--squash", branchName)
	exit.IfWrap(err, "Error squash merging")
}
