package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
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

func NewBranchType(existingOpt Option[configdomain.NewBranchType], inputs dialogcomponents.TestInputs) (Option[configdomain.NewBranchType], dialogdomain.Exit, error) {
	entries := list.Entries[Option[configdomain.BranchType]]{
		{
			Data: Some(configdomain.BranchTypeFeatureBranch),
			Text: "always create feature branches",
		},
		{
			Data: Some(configdomain.BranchTypeParkedBranch),
			Text: "always create parked branches",
		},
		{
			Data: Some(configdomain.BranchTypePrototypeBranch),
			Text: "always create prototype branches",
		},
		{
			Data: Some(configdomain.BranchTypePerennialBranch),
			Text: "always create perennial branches",
		},
	}
	existingOptBranchType := None[configdomain.BranchType]()
	if existing, has := existingOpt.Get(); has {
		existingOptBranchType = Some(existing.BranchType())
	}
	defaultPos := entries.IndexOf(existingOptBranchType)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, newBranchTypeTitle, NewBranchTypeHelp, inputs, "new-branch-type")
	fmt.Println(messages.NewBranchType, dialogcomponents.FormattedSelection(selection.String(), exit))
	return configdomain.NewBranchTypeOption(selection), exit, err
}
