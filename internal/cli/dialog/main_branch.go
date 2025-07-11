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
func MainBranch(args MainBranchArgs) (Option[gitdomain.LocalBranchName], dialogdomain.Exit, error) {
	// if set in global config: add option "use global setting" with None
	// if set in local config: don't add None option, preselect local setting
	// if no local config: don't add None option, keep existing preselect nothing
	cursor := 0
	if defaultEntry, hasDefaultEntry := args.DefaultEntryOpt.Get(); hasDefaultEntry {
		cursor = slice.Index(args.LocalBranches, defaultEntry).GetOrElse(0)
	}
	selection, exit, err := dialogcomponents.RadioList(list.NewEntries(args.LocalBranches...), cursor, mainBranchTitle, MainBranchHelp, args.Inputs)
	fmt.Printf(messages.MainBranch, dialogcomponents.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}

type MainBranchArgs struct {
	GitStandardBranch   Option[gitdomain.LocalBranchName]
	GlobalGitMainBranch Option[gitdomain.LocalBranchName]
	LocalGitMainBranch  Option[gitdomain.LocalBranchName]
	LocalBranches       gitdomain.LocalBranchNames
	Inputs              dialogcomponents.TestInput
}
