package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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
func MainBranch(args MainBranchArgs) (MainBranchResult, dialogdomain.Exit, error) {
	// populate the local branches
	entries := list.Entries[Option[gitdomain.LocalBranchName]]{}
	unscoped, hasUnscoped := args.Unscoped.Get()
	local, hasLocal := args.Local.Get()
	if hasUnscoped && !hasLocal {
		entries = append(entries, list.Entry[Option[gitdomain.LocalBranchName]]{
			Data: None[gitdomain.LocalBranchName](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, unscoped),
		})
	}
	for _, localBranch := range args.LocalBranches {
		entries = append(entries, list.Entry[Option[gitdomain.LocalBranchName]]{
			Data: Some(localBranch),
			Text: localBranch.String(),
		})
	}

	// pre-select the already configured value
	var cursor int
	switch {
	case hasLocal:
		cursor = entries.IndexOf(Some(local))
	case hasUnscoped:
		cursor = 0
	default:
		cursor = entries.IndexOf(args.StandardBranch)
	}

	// show the dialog
	selection, exit, err := dialogcomponents.RadioList(entries, cursor, mainBranchTitle, MainBranchHelp, args.Inputs, "main-branch")
	if err == nil {
		fmt.Printf(messages.MainBranch, dialogcomponents.FormattedOption(selection, hasUnscoped, exit))
	}
	return MainBranchResult{
		ActualMainBranch: selection.GetOr(unscoped), // the user either selected a branch, or None if unscoped exists
		UserChoice:       selection,
	}, exit, err
}

type MainBranchArgs struct {
	Inputs         dialogcomponents.Inputs
	Local          Option[gitdomain.LocalBranchName]
	LocalBranches  gitdomain.LocalBranchNames
	StandardBranch Option[gitdomain.LocalBranchName]
	Unscoped       Option[gitdomain.LocalBranchName]
}

type MainBranchResult struct {
	ActualMainBranch gitdomain.LocalBranchName
	UserChoice       Option[gitdomain.LocalBranchName]
}
