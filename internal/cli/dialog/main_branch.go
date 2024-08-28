package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v15/internal/cli/dialog/components"
	"github.com/git-town/git-town/v15/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/messages"
	. "github.com/git-town/git-town/v15/pkg/prelude"
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
func MainBranch(localBranches gitdomain.LocalBranchNames, defaultEntryOpt Option[gitdomain.LocalBranchName], inputs components.TestInput) (gitdomain.LocalBranchName, bool, error) {
	cursor := 0
	if defaultEntry, hasDefaultEntry := defaultEntryOpt.Get(); hasDefaultEntry {
		cursor = localBranches.IndexOr(defaultEntry, 0)
	}
	selection, aborted, err := components.RadioList(list.NewEntries(localBranches...), cursor, mainBranchTitle, MainBranchHelp, inputs)
	fmt.Printf(messages.MainBranch, components.FormattedSelection(selection.String(), aborted))
	return selection, aborted, err
}
