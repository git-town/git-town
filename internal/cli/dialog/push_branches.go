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
	pushBranchesTitle = `Push branches`
	PushBranchesHelp  = `
The push-branches setting controls whether
Git Town pushes branches to the development remote
while syncing.

Pushing is enabled by default.
More details: https://www.git-town.com/preferences/push-branches.

`
)

func PushBranches(args Args[configdomain.PushBranches]) (Option[configdomain.PushBranches], dialogdomain.Exit, error) {
	entries := list.Entries[Option[configdomain.PushBranches]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[configdomain.PushBranches]]{
			Data: None[configdomain.PushBranches](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[configdomain.PushBranches]]{
		{
			Data: Some(configdomain.PushBranches(true)),
			Text: "enabled: push branches while syncing",
		},
		{
			Data: Some(configdomain.PushBranches(false)),
			Text: "disabled: don't push branches while syncing",
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, pushBranchesTitle, PushBranchesHelp, args.Inputs, "push-branches")
	fmt.Printf(messages.PushBranches, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
