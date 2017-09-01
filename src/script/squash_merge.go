package script

import "github.com/Originate/git-town/src/logs"

// SquashMerge squash merges the given branch into the current branch
func SquashMerge(branchName string) {
	err := RunCommand("git", "merge", "--squash", branchName)
	logs.FatalOnWrap(err, "Error squash merging")
}
