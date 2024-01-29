package enter

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/dialogs/components"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/gohacks/stringers"
)

const enterBranchHelp = `
The main branch is the branch from which you cut new feature branches,
and into which you ship feature branches when they are done.
In most repositories, this branch is called "main", "master", or "development".


`

// MainBranch lets the user select a new main branch for this repo.
func MainBranch(localBranches gitdomain.LocalBranchNames, oldMainBranch gitdomain.LocalBranchName, inputs components.TestInput) (gitdomain.LocalBranchName, bool, error) {
	cursor := stringers.IndexOrStart(localBranches, oldMainBranch)
	selection, aborted, err := components.RadioList(localBranches, cursor, enterBranchHelp, inputs)
	fmt.Printf("Main branch: %s\n", components.FormattedSelection(selection.String(), aborted))
	return selection, aborted, err
}
