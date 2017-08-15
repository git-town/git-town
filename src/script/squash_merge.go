package script

import "log"

// SquashMerge squash merges the given branch into the current branch
func SquashMerge(branchName string) {
	err := RunCommand("git", "merge", "--squash", branchName)
	if err != nil {
		log.Fatal("Error squash merging:", err)
	}
}
