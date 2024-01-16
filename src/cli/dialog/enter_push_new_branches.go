package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

const enterPushNewBranchesHelp = `
After creating new feature branches
by running "git hack", "git append", or "git prepend",
should Git Town push the new branch to the origin remote?
Doing so will trigger an unnecessary CI run that tests the empty branch.
The extra network access also increases the time these commands take.


`

// EnterMainBranch lets the user select a new main branch for this repo.
func EnterPushNewBranches(localBranches gitdomain.LocalBranchNames, oldMainBranch gitdomain.LocalBranchName, inputs TestInput) (gitdomain.LocalBranchName, bool, error) {
	selection, aborted, err := radioList(radioListArgs{
		entries:      localBranches.Strings(),
		defaultEntry: oldMainBranch.String(),
		help:         enterBranchHelp,
		testInput:    inputs,
	})
	fmt.Printf("Selected main branch: %s\n", formattedSelection(selection, aborted))
	return gitdomain.LocalBranchName(selection), aborted, err
}
