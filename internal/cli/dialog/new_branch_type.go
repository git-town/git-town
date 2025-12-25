package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

const (
	newBranchTypeTitle = `New branch type`
	NewBranchTypeHelp  = `
This setting controls the type of new branches
that you create with "git town hack", "append", or "prepend".
Defaults to feature branches.

More details: https://www.git-town.com/preferences/new-branch-type.

`
)

func NewBranchType(args Args[configdomain.NewBranchType]) (Option[configdomain.NewBranchType], dialogdomain.Exit, error) {
	entries := list.Entries[Option[configdomain.BranchType]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[configdomain.BranchType]]{
			Data: None[configdomain.BranchType](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[configdomain.BranchType]]{
		{
			Data: Some(configdomain.BranchTypeFeatureBranch),
			Text: "create feature branches",
		},
		{
			Data: Some(configdomain.BranchTypeParkedBranch),
			Text: "create parked branches",
		},
		{
			Data: Some(configdomain.BranchTypePrototypeBranch),
			Text: "create prototype branches",
		},
		{
			Data: Some(configdomain.BranchTypePerennialBranch),
			Text: "create perennial branches",
		},
	}...)
	cursor := 0
	if local, hasLocal := args.Local.Get(); hasLocal {
		cursor = entries.IndexOf(Some(local.BranchType()))
	}
	input, exit, err := dialogcomponents.RadioList(entries, cursor, newBranchTypeTitle, NewBranchTypeHelp, args.Inputs, "new-branch-type")
	selection := configdomain.NewBranchTypeOpt(input)
	fmt.Println(messages.NewBranchType, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
