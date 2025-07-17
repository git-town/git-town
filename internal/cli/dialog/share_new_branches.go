package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
)

const (
	shareNewBranchesTitle = `Share new branches`
	ShareNewBranchesHelp  = `
How should Git Town share the new branches it creates?

Possible options:

- none: New branches are local until you sync or propose them
- push: push new branches to the development remote
- propose: propose new branches

`
)

func ShareNewBranches(existing configdomain.ShareNewBranches, inputs dialogcomponents.TestInputs) (configdomain.ShareNewBranches, dialogdomain.Exit, error) {
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
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, shareNewBranchesTitle, ShareNewBranchesHelp, inputs, "share-new-branches")
	fmt.Printf(messages.ShareNewBranches, dialogcomponents.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}
