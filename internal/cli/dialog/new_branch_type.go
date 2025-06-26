package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components/list"
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

func NewBranchType(existingOpt Option[configdomain.BranchType], inputs components.TestInput) (Option[configdomain.BranchType], dialogdomain.Exit, error) {
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
	defaultPos := 0
	for e, entry := range entries {
		if entry.Data.Equal(existingOpt) {
			defaultPos = e
		}
	}
	selection, exit, err := components.RadioList(entries, defaultPos, newBranchTypeTitle, NewBranchTypeHelp, inputs)
	if err != nil || exit {
		return None[configdomain.BranchType](), exit, err
	}
	fmt.Println(messages.CreatePrototypeBranches, components.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}
