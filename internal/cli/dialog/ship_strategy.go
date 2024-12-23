package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/messages"
)

const (
	shipStrategyTitle = `Ship strategy`
	ShipStrategyHelp  = `
Which method should Git Town use to ship feature branches?

Options:

- api: merge the proposal on your code hosting platform via the code hosting API
- fast-forward: in your local repo, fast-forward the parent branch to point to the commits on the feature branch
- squash-merge: in your local repo, squash-merge the feature branch into its parent branch

All options update proposals of child branches and remove the shipped branch locally and remotely.
`
)

func ShipStrategy(existing configdomain.ShipStrategy, inputs components.TestInput) (configdomain.ShipStrategy, bool, error) {
	entries := list.Entries[configdomain.ShipStrategy]{
		{
			Data: configdomain.ShipStrategyAPI,
			Text: `api: merge the proposal on your code hosting platform via the code hosting API`,
		},
		{
			Data: configdomain.ShipStrategyFastForward,
			Text: `fast-forward: in your local repo, fast-forward the parent branch to point to the commits on the feature branch`,
		},
		{
			Data: configdomain.ShipStrategySquashMerge,
			Text: `squash-merge: in your local repo, squash-merge the feature branch into its parent branch`,
		},
	}
	defaultPos := shipStrategyEntryIndex(entries, existing)
	selection, aborted, err := components.RadioList(entries, defaultPos, shipStrategyTitle, ShipStrategyHelp, inputs)
	if err != nil || aborted {
		return configdomain.ShipStrategyAPI, aborted, err
	}
	fmt.Printf(messages.ShipDeletesTrackingBranches, components.FormattedSelection(selection.String(), aborted))
	return selection, aborted, err
}

func shipStrategyEntryIndex(entries list.Entries[configdomain.ShipStrategy], needle configdomain.ShipStrategy) int {
	for e, entry := range entries {
		if entry.Data == needle {
			return e
		}
	}
	return 0
}
