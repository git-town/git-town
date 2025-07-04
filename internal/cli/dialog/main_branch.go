package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/slice"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
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
func MainBranch(localBranches gitdomain.LocalBranchNames, defaultEntryOpt Option[gitdomain.LocalBranchName], inputs dialogcomponents.TestInput) (gitdomain.LocalBranchName, dialogdomain.Exit, error) {
	cursor := 0
	if defaultEntry, hasDefaultEntry := defaultEntryOpt.Get(); hasDefaultEntry {
		cursor = slice.Index(localBranches, defaultEntry).GetOrElse(0)
	}
	selection, exit, err := dialogcomponents.RadioList(list.NewEntries(localBranches...), cursor, mainBranchTitle, MainBranchHelp, inputs)
	fmt.Printf(messages.MainBranch, dialogcomponents.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}
