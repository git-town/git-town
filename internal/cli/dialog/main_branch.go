package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v20/internal/cli/dialog/components"
	"github.com/git-town/git-town/v20/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/gohacks/slice"
	"github.com/git-town/git-town/v20/internal/messages"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

const (
	mainBranchTitle = `Main branch`
	MainBranchHelp  = `
The main branch is your project's default branch.
It's where new feature branches are created from,
and where completed features are merged back into.

This is typically the branch called
"main", "master", or "development".

`
)

// MainBranch lets the user select a new main branch for this repo.
func MainBranch(localBranches gitdomain.LocalBranchNames, defaultEntryOpt Option[gitdomain.LocalBranchName], inputs components.TestInput) (gitdomain.LocalBranchName, bool, error) {
	cursor := 0
	if defaultEntry, hasDefaultEntry := defaultEntryOpt.Get(); hasDefaultEntry {
		cursor = slice.Index(localBranches, defaultEntry).GetOrElse(0)
	}
	selection, aborted, err := components.RadioList(list.NewEntries(localBranches...), cursor, mainBranchTitle, MainBranchHelp, inputs)
	fmt.Printf(messages.MainBranch, components.FormattedSelection(selection.String(), aborted))
	return selection, aborted, err
}
