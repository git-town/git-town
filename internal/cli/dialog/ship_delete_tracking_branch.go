package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/cli/format"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
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

func ShipDeleteTrackingBranch(existing configdomain.ShipDeleteTrackingBranch, inputs dialogcomponents.TestInput) (configdomain.ShipDeleteTrackingBranch, dialogdomain.Exit, error) {
	entries := list.Entries[bool]{
		{
			Data: true,
			Text: `yes, "git town ship" should delete tracking branches`,
		},
		{
			Data: false,
			Text: `no, my forge deletes branches after merging them`,
		},
	}
	var defaultPos int
	if existing {
		defaultPos = 0
	} else {
		defaultPos = 1
	}
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, shipDeleteTrackingBranchTitle, ShipDeleteTrackingBranchHelp, inputs)
	fmt.Printf(messages.ShipDeletesTrackingBranches, dialogcomponents.FormattedSelection(format.Bool(selection), exit))
	return configdomain.ShipDeleteTrackingBranch(selection), exit, err
}
