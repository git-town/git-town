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
	shipDeleteTrackingBranchTitle = `Ship delete tracking branch`
	ShipDeleteTrackingBranchHelp  = `
Should git town ship delete
the remote tracking branch after shipping?

Disable this if your code hosting provider
(GitHub, GitLab, etc.) automatically deletes
branches when pull requests are merged through its UI.

`
)

func ShipDeleteTrackingBranch(args Args[configdomain.ShipDeleteTrackingBranch]) (Option[configdomain.ShipDeleteTrackingBranch], dialogdomain.Exit, error) {
	entries := list.Entries[Option[configdomain.ShipDeleteTrackingBranch]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[configdomain.ShipDeleteTrackingBranch]]{
			Data: None[configdomain.ShipDeleteTrackingBranch](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[configdomain.ShipDeleteTrackingBranch]]{
		{
			Data: Some(configdomain.ShipDeleteTrackingBranch(true)),
			Text: `yes, "git town ship" should delete tracking branches`,
		},
		{
			Data: Some(configdomain.ShipDeleteTrackingBranch(false)),
			Text: `no, my forge deletes branches after merging them`,
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, shipDeleteTrackingBranchTitle, ShipDeleteTrackingBranchHelp, args.Inputs, "ship-delete-tracking-branch")
	fmt.Printf(messages.ShipDeletesTrackingBranches, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
