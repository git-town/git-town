package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/gohacks/stringers"
)

const enterBranchHelp = `
The main branch is the branch from which you cut new feature branches,
and into which you ship feature branches when they are done.
In most repositories, this branch is called "main", "master", or "development".


`

// EnterMainBranch lets the user select a new main branch for this repo.
func EnterMainBranch(localBranches gitdomain.LocalBranchNames, oldMainBranch gitdomain.LocalBranchName, inputs TestInput) (gitdomain.LocalBranchName, bool, error) {
	cursor := stringers.IndexOrStart(localBranches, oldMainBranch)
	selection, aborted, err := radioList(localBranches, cursor, enterBranchHelp, inputs)
	fmt.Printf("Main branch: %s\n", formattedSelection(selection.String(), aborted))
	return selection, aborted, err
}
