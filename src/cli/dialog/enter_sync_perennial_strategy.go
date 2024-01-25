package dialog

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const enterSyncPerennialStrategyHelp = `
How should Git Town synchronize perennial branches?

Perennial branches don't have a parent branch.
The only updates they receive are additional commits
to their tracking branch made somewhere else.

`

const (
	SyncPerennialStrategyEntryMerge  syncPerennialStrategyEntry = `merge updates from the tracking branch into perennial branches`
	SyncPerennialStrategyEntryRebase syncPerennialStrategyEntry = `rebase perennial branches against their tracking branch`
)

func EnterSyncPerennialStrategy(existing configdomain.SyncPerennialStrategy, inputs TestInput) (configdomain.SyncPerennialStrategy, bool, error) {
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
	selection, aborted, err := radioList(entries, defaultPos, enterSyncPerennialStrategyHelp, inputs)
	if err != nil || aborted {
		return configdomain.SyncPerennialStrategyRebase, aborted, err
	}
	fmt.Printf("Sync perennial branches: %s\n", formattedSelection(selection.Short(), aborted))
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
