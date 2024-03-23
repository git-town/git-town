package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v13/src/cli/dialog/components"
	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/src/gohacks/stringers"
	"github.com/git-town/git-town/v13/src/messages"
)

const (
	mainBranchTitle = `Main branch`
	MainBranchHelp  = `
The main branch is the branch from which you cut new feature branches,
and into which you ship feature branches when they are done.
This branch is often called "main", "master", or "development".

`
)

// MainBranch lets the user select a new main branch for this repo.
func MainBranch(localBranches gitdomain.LocalBranchNames, defaultEntry gitdomain.LocalBranchName, inputs components.TestInput) (gitdomain.LocalBranchName, bool, error) {
	cursor := stringers.IndexOrStart(localBranches, defaultEntry)
	selection, aborted, err := components.RadioList(localBranches, cursor, mainBranchTitle, MainBranchHelp, inputs)
	fmt.Printf(messages.MainBranch, components.FormattedSelection(selection.String(), aborted))
	return selection, aborted, err
}
