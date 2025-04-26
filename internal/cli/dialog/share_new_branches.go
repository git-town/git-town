package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v19/internal/cli/dialog/components"
	"github.com/git-town/git-town/v19/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/messages"
)

const (
	pushNewBranchesTitle = `Share new branches`
	PushNewBranchesHelp  = `
How should Git Town share the new branches it creates?

Possible options:

- none: New branches remain local until you sync or propose them.
- push: New branches are automatically pushed to the development remote.

`
)

func PushNewBranches(existing configdomain.ShareNewBranches, inputs components.TestInput) (configdomain.ShareNewBranches, bool, error) {
	entries := list.Entries[configdomain.ShareNewBranches]{
		{
			Data: configdomain.ShareNewBranchesNone,
			Text: "yes: push new branches to origin",
		},
		{
			Data: configdomain.ShareNewBranchesPush,
			Text: "no, new branches remain local until synced",
		},
	}
	defaultPos := entries.IndexOf(existing)
	selection, aborted, err := components.RadioList(entries, defaultPos, pushNewBranchesTitle, PushNewBranchesHelp, inputs)
	if err != nil || aborted {
		return configdomain.ShareNewBranchesNone, aborted, err
	}
	fmt.Printf(messages.PushNewBranches, components.FormattedSelection(selection.String(), aborted))
	return selection, aborted, err
}
