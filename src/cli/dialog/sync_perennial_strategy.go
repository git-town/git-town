package dialog

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/messages"
)

const (
	syncPerennialStrategyTitle = `Sync-perennial strategy`
	SyncPerennialStrategyHelp  = `
How should Git Town synchronize perennial branches?
Perennial branches have no parent branch.
The only updates they receive are additional commits
made to their tracking branch somewhere else.

`
)

const (
	SyncPerennialStrategyEntryMerge  syncPerennialStrategyEntry = `merge updates from the tracking branch into perennial branches`
	SyncPerennialStrategyEntryRebase syncPerennialStrategyEntry = `rebase perennial branches against their tracking branch`
)

func SyncPerennialStrategy(existing configdomain.SyncPerennialStrategy, inputs components.TestInput) (configdomain.SyncPerennialStrategy, bool, error) {
	entries := []syncPerennialStrategyEntry{
		SyncPerennialStrategyEntryMerge,
		SyncPerennialStrategyEntryRebase,
	}
	var defaultPos int
	switch existing {
	case configdomain.SyncPerennialStrategyMerge:
		defaultPos = 0
	case configdomain.SyncPerennialStrategyRebase:
		defaultPos = 1
	default:
		panic("unknown sync-perennial-strategy: " + existing.String())
	}
	selection, aborted, err := components.RadioList(components.NewEnabledBubbleListEntries(entries), defaultPos, syncPerennialStrategyTitle, SyncPerennialStrategyHelp, inputs)
	if err != nil || aborted {
		return configdomain.SyncPerennialStrategyRebase, aborted, err
	}
	fmt.Printf(messages.SyncPerennialBranches, components.FormattedSelection(selection.Short(), aborted))
	return selection.SyncPerennialStrategy(), aborted, err
}

type syncPerennialStrategyEntry string

func (self syncPerennialStrategyEntry) Short() string {
	start, _, _ := strings.Cut(self.String(), " ")
	return start
}

func (self syncPerennialStrategyEntry) String() string {
	return string(self)
}

func (self syncPerennialStrategyEntry) SyncPerennialStrategy() configdomain.SyncPerennialStrategy {
	switch self {
	case SyncPerennialStrategyEntryMerge:
		return configdomain.SyncPerennialStrategyMerge
	case SyncPerennialStrategyEntryRebase:
		return configdomain.SyncPerennialStrategyRebase
	}
	panic("unhandled syncPerennialStrategyEntry: " + self)
}
