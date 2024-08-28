package dialog

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v15/internal/cli/dialog/components"
	"github.com/git-town/git-town/v15/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/messages"
)

const (
	shipStrategyTitle = `Ship strategy`
	ShipStrategyHelp  = `
Which method should Git Town use to ship feature branches?

Options:

- api: Git Town presses the "merge" button on your code hosting platform for you by talking to the code hosting API
- squash-merge: Git Town squash-merges the feature branch into its parent branch on your local machine

All options update proposals of child branches and remove the shipped branch locally and remotely.
`
)

const (
	ShipStrategyEntryAPI         shipStrategyEntry = `api: "git ship" presses the "merge" button on your code hosting platform for you`
	ShipStrategyEntrySquashMerge shipStrategyEntry = `squash-merge: "git ship" squash-merges the branch on your local machine`
)

func ShipStrategy(existing configdomain.ShipStrategy, inputs components.TestInput) (configdomain.ShipStrategy, bool, error) {
	entries := []shipStrategyEntry{
		ShipStrategyEntryAPI,
		ShipStrategyEntrySquashMerge,
	}
	defaultPos := shipStrategyEntryIndex(entries, existing)
	selection, aborted, err := components.RadioList(list.NewEntries(entries...), defaultPos, shipStrategyTitle, ShipStrategyHelp, inputs)
	if err != nil || aborted {
		return configdomain.ShipStrategyAPI, aborted, err
	}
	fmt.Printf(messages.ShipDeletesTrackingBranches, components.FormattedSelection(selection.Short(), aborted))
	return selection.ShipStrategy(), aborted, err
}

type shipStrategyEntry string

func (self shipStrategyEntry) ShipStrategy() configdomain.ShipStrategy {
	switch self {
	case ShipStrategyEntryAPI:
		return configdomain.ShipStrategyAPI
	case ShipStrategyEntrySquashMerge:
		return configdomain.ShipStrategySquashMerge
	}
	panic("unhandled shipStrategyEntry: " + self)
}

func (self shipStrategyEntry) Short() string {
	start, _, _ := strings.Cut(self.String(), ":")
	return start
}

func (self shipStrategyEntry) String() string {
	return string(self)
}

func shipStrategyEntryIndex(entries []shipStrategyEntry, needle configdomain.ShipStrategy) int {
	needleText := needle.String()
	for e, entry := range entries {
		if entry.Short() == needleText {
			return e
		}
	}
	return 0
}
