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
	shareNewBranchesTitle = `Share new branches`
	ShareNewBranchesHelp  = `
How should Git Town share the new branches it creates?

Possible options:

- none: New branches are local until you sync or propose them
- push: push new branches to the development remote
- propose: propose new branches

`
)

func ShareNewBranches(args Args[configdomain.ShareNewBranches]) (Option[configdomain.ShareNewBranches], dialogdomain.Exit, error) {
	entries := list.Entries[Option[configdomain.ShareNewBranches]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[configdomain.ShareNewBranches]]{
			Data: None[configdomain.ShareNewBranches](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[configdomain.ShareNewBranches]]{
		{
			Data: Some(configdomain.ShareNewBranchesNone),
			Text: "no sharing: new branches remain local until synced or proposed",
		},
		{
			Data: Some(configdomain.ShareNewBranchesPush),
			Text: "push new branches to the dev remote",
		},
		{
			Data: Some(configdomain.ShareNewBranchesPropose),
			Text: "propose new branches",
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, shareNewBranchesTitle, ShareNewBranchesHelp, args.Inputs, "share-new-branches")
	fmt.Printf(messages.ShareNewBranches, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
