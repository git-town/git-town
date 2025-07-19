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
	// if local: don't add None option, preselect local setting
	// if no local and unscoped: add None option and preselect it
	// if local and different unscoped: add None option but preselect the local one
	entries := list.Entries[Option[gitdomain.LocalBranchName]]{}
	if unscopedMain, hasUnscoped := args.UnscopedGitMainBranch.Get(); hasUnscoped {
		if args.LocalGitMainBranch.IsNone() {
			if len(args.LocalBranches) == 1 && unscopedMain == args.LocalBranches[0] {
				return None[gitdomain.LocalBranchName](), false, nil
			}
			entries = append(entries, list.Entry[Option[gitdomain.LocalBranchName]]{
				Data: None[gitdomain.LocalBranchName](),
				Text: fmt.Sprintf("use global setting (%s)", unscopedMain),
			})
		}
	}
	for _, localBranch := range args.LocalBranches {
		entries = append(entries, list.Entry[Option[gitdomain.LocalBranchName]]{
			Data: Some(localBranch),
			Text: localBranch.String(),
		})
	}
	cursor := 0
	if gitStandard, hasStandard := args.GitStandardBranch.Get(); hasStandard {
		if index, hasIndex := slice.Index(args.LocalBranches, gitStandard).Get(); hasIndex {
			cursor = index
		}
	}
	if localMain, hasLocal := args.LocalGitMainBranch.Get(); hasLocal {
		if index, hasIndex := slice.Index(args.LocalBranches, localMain).Get(); hasIndex {
			cursor = index
		}
	}
	selection, exit, err := dialogcomponents.RadioList(entries, cursor, mainBranchTitle, MainBranchHelp, args.Inputs, "main-branch")
	fmt.Printf(messages.MainBranch, dialogcomponents.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}

type MainBranchArgs struct {
	GitStandardBranch     Option[gitdomain.LocalBranchName]
	Inputs                dialogcomponents.TestInputs
	LocalBranches         gitdomain.LocalBranchNames
	LocalGitMainBranch    Option[gitdomain.LocalBranchName]
	UnscopedGitMainBranch Option[gitdomain.LocalBranchName]
}
