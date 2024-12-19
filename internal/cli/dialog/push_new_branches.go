package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v17/internal/cli/format"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/messages"
)

const (
	pushNewBranchesTitle = `Push new branches`
	PushNewBranchesHelp  = `
Should Git Town push the new branches it creates
immediately to origin even if they are empty?

When enabled, you can run "git push" right away
but creating new branches is slower and
it triggers an unnecessary CI run on the empty branch.

When disabled, many Git Town commands execute faster
and Git Town will create the missing tracking branch
on the first run of "git town sync".

`
)

func PushNewBranches(existing configdomain.PushNewBranches, inputs components.TestInput) (configdomain.PushNewBranches, bool, error) {
	entries := list.Entries[configdomain.PushNewBranches]{
		{
			Data: true,
			Text: "yes: push new branches to origin",
		},
		{
			Data: false,
			Text: "no, new branches remain local until synced",
		},
	}
	defaultPos := entries.IndexOf(existing)
	selection, aborted, err := components.RadioList(entries, defaultPos, pushNewBranchesTitle, PushNewBranchesHelp, inputs)
	if err != nil || aborted {
		return true, aborted, err
	}
	fmt.Printf(messages.PushNewBranches, components.FormattedSelection(format.Bool(selection.IsTrue()), aborted))
	return selection, aborted, err
}
