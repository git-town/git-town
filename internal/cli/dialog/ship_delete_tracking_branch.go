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
	shipDeleteTrackingBranchTitle = `Ship delete tracking branch`
	ShipDeleteTrackingBranchHelp  = `
Should git town ship delete the remote tracking branch after shipping?

Disable this if your code hosting provider (GitHub, GitLab, etc.)
automatically deletes branches when pull requests are merged through its UI.

`
)

func ShipDeleteTrackingBranch(existing configdomain.ShipDeleteTrackingBranch, inputs components.TestInput) (configdomain.ShipDeleteTrackingBranch, bool, error) {
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
	selection, aborted, err := components.RadioList(entries, defaultPos, shipDeleteTrackingBranchTitle, ShipDeleteTrackingBranchHelp, inputs)
	if err != nil || aborted {
		return true, aborted, err
	}
	fmt.Printf(messages.ShipDeletesTrackingBranches, components.FormattedSelection(format.Bool(selection), aborted))
	return configdomain.ShipDeleteTrackingBranch(selection), aborted, err
}
