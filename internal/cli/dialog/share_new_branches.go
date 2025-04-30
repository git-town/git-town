package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v19/internal/cli/dialog/components"
	"github.com/git-town/git-town/v19/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/messages"
)

const (
	shareNewBranchesTitle = `Share new branches`
	ShareNewBranchesHelp  = `
How should Git Town share the new branches it creates?

Possible options:

- none: New branches remain local until you sync or propose them.
- push: push new branches to the development remote.
- propose: propose new branches

`
)

func ShareNewBranches(existing configdomain.ShareNewBranches, inputs components.TestInput) (configdomain.ShareNewBranches, bool, error) {
	entries := list.Entries[configdomain.ShareNewBranches]{
		{
			Data: configdomain.ShareNewBranchesNone,
			Text: "no sharing: new branches remain local until synced or proposed",
		},
		{
			Data: configdomain.ShareNewBranchesPush,
			Text: "push new branches to the dev remote",
		},
		{
			Data: configdomain.ShareNewBranchesPropose,
			Text: "propose new branches",
		},
	}
	defaultPos := entries.IndexOf(existing)
	selection, aborted, err := components.RadioList(entries, defaultPos, shareNewBranchesTitle, ShareNewBranchesHelp, inputs)
	if err != nil || aborted {
		return configdomain.ShareNewBranchesNone, aborted, err
	}
	fmt.Printf(messages.ShareNewBranches, components.FormattedSelection(selection.String(), aborted))
	return selection, aborted, err
}
