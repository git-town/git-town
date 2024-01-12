package dialog

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

const enterBranchHelp = `
Let's start by configuring the main development branch.
This is the branch from which you cut new feature branches,
and into which you ship feature branches when they are done.
In most repositories, this is the "main", "master", or "development" branch.


`

// EnterMainBranch lets the user select a new main branch for this repo.
func EnterMainBranch(localBranches gitdomain.LocalBranchNames, oldMainBranch gitdomain.LocalBranchName, inputs TestInput) (gitdomain.LocalBranchName, bool, error) {
	selection, aborted, err := radioList(radioListArgs{
		entries:      localBranches.Strings(),
		defaultEntry: oldMainBranch.String(),
		help:         enterBranchHelp,
		testInput:    inputs,
	})
	return gitdomain.LocalBranchName(selection), aborted, err
}
