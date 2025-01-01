package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/messages"
	. "github.com/git-town/git-town/v17/pkg/prelude"
)

const (
	newBranchTypeTitle = `New branch type`
	NewBranchTypeHelp  = `
The "new-branch-type" setting allows you to override the type that new branches
will have when you run "git town hack", "append", or "prepend".

Branches for which no branch type is set, and for which no configuration entries match, are considered feature branches.
More info at https://www.git-town.com/preferences/new-branch-type.

`
)

func NewBranchType(existingOpt Option[configdomain.BranchType], inputs components.TestInput) (Option[configdomain.BranchType], bool, error) {
	entries := list.Entries[Option[configdomain.BranchType]]{
		{
			Data: None[configdomain.BranchType](),
			Text: "create default branch type",
		},
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
			Data: Some(configdomain.BranchTypePrototypeBranch),
			Text: "always create perennial branches",
		},
	}
	defaultPos := entries.IndexOf(existingOpt)
	selection, aborted, err := components.RadioList(entries, defaultPos, newBranchTypeTitle, NewBranchTypeHelp, inputs)
	if err != nil || aborted {
		return None[configdomain.BranchType](), aborted, err
	}
	fmt.Println(messages.CreatePrototypeBranches, components.FormattedSelection(selection.String(), aborted))
	return selection, aborted, err
}
