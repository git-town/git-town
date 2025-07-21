package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
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
func MainBranch(args MainBranchArgs) (selectedMainBranch Option[gitdomain.LocalBranchName], mainBranch gitdomain.LocalBranchName, exit dialogdomain.Exit, err error) {
	// populate the local branches
	entries := make(list.Entries[Option[gitdomain.LocalBranchName]], 0, len(args.LocalBranches)+1)
	for _, localBranch := range args.LocalBranches {
		entries = append(entries, list.Entry[Option[gitdomain.LocalBranchName]]{
			Data: Some(localBranch),
			Text: localBranch.String(),
		})
	}

	// optionally add "None" entry and pre-select the already configured value
	cursor := 0
	unscopedMain, hasUnscoped := args.UnscopedGitMainBranch.Get()
	local, hasLocal := args.LocalGitMainBranch.Get()
	switch {
	case !hasLocal && !hasUnscoped:
		cursor = entries.IndexOf(args.GitStandardBranch)
	case hasLocal && !hasUnscoped:
		cursor = entries.IndexOf(Some(local))
	case !hasLocal && hasUnscoped:
		noneEntry := list.Entry[Option[gitdomain.LocalBranchName]]{
			Data: None[gitdomain.LocalBranchName](),
			Text: fmt.Sprintf(messages.DialogResultUseGlobalValue, unscopedMain),
		}
		entries = append(list.Entries[Option[gitdomain.LocalBranchName]]{noneEntry}, entries...)
		cursor = 0
	case hasLocal && hasUnscoped:
		noneEntry := list.Entry[Option[gitdomain.LocalBranchName]]{
			Data: None[gitdomain.LocalBranchName](),
			Text: fmt.Sprintf(messages.DialogResultUseGlobalValue, unscopedMain),
		}
		entries = append(list.Entries[Option[gitdomain.LocalBranchName]]{noneEntry}, entries...)
		cursor = entries.IndexOf(Some(local))
	}

	// show the dialog
	selection, exit, err := dialogcomponents.RadioList(entries, cursor, mainBranchTitle, MainBranchHelp, args.Inputs, "main-branch")
	fmt.Printf(messages.MainBranch, dialogcomponents.FormattedSelection(selection.String(), exit))
	mainBranch = selection.GetOrElse(unscopedMain) // the user either selected a branch, or None if unscoped exists
	return selection, mainBranch, exit, err
}

type MainBranchArgs struct {
	GitStandardBranch     Option[gitdomain.LocalBranchName]
	Inputs                dialogcomponents.TestInputs
	LocalBranches         gitdomain.LocalBranchNames
	LocalGitMainBranch    Option[gitdomain.LocalBranchName]
	UnscopedGitMainBranch Option[gitdomain.LocalBranchName]
}
