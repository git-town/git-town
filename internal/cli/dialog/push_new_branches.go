package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v19/internal/cli/dialog/components"
	"github.com/git-town/git-town/v19/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v19/internal/cli/format"
	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/messages"
)

const (
	pushNewBranchesTitle = `Sharing new branches`
	PushNewBranchesHelp  = `
How should Git Town share the new branches it creates?

Possible options:

- none: New branches remain local until you sync or propose them.
- push: New branches are automatically pushed to the development remote.

`
)

func PushNewBranches(existing configdomain.PushNewBranches, inputs components.TestInput) (configdomain.PushNewBranches, bool, error) {
	entries := list.Entries[configdomain.PushNewBranches]{
		{
			Data: true,
			Text: "don't share: new branches are local until synced or proposed",
		},
		{
			Data: false,
			Text: "push new branches to the dev remote",
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
