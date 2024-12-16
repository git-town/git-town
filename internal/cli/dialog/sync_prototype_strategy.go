package dialog

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/messages"
)

const (
	syncPrototypeStrategyTitle = `Sync-prototype strategy`
	SyncPrototypeStrategyHelp  = `
How should Git Town synchronize prototype branches?
Prototype branches are feature branches that haven't been proposed yet.
Typically they contain  features and bug fixes on them,
hence their name.

`
)

const (
	syncPrototypeStrategyEntryMerge    syncPrototypeStrategyEntry = `merge updates from the parent and tracking branch`
	syncPrototypeStrategyEntryRebase   syncPrototypeStrategyEntry = `rebase branches against their parent and tracking branch`
	syncPrototypeStrategyEntryCompress syncPrototypeStrategyEntry = `compress the branch after merging parent and tracking`
)

func SyncPrototypeStrategy(existing configdomain.SyncPrototypeStrategy, inputs components.TestInput) (configdomain.SyncPrototypeStrategy, bool, error) {
	entries := []syncPrototypeStrategyEntry{
		syncPrototypeStrategyEntryMerge,
		syncPrototypeStrategyEntryRebase,
		syncPrototypeStrategyEntryCompress,
	}
	var defaultPos int
	switch existing {
	case configdomain.SyncPrototypeStrategyMerge:
		defaultPos = 0
	case configdomain.SyncPrototypeStrategyRebase:
		defaultPos = 1
	case configdomain.SyncPrototypeStrategyCompress:
		defaultPos = 2
	default:
		panic("unknown sync strategy: " + existing.String())
	}
	selection, aborted, err := components.RadioList(list.NewEntries(entries...), defaultPos, syncPrototypeStrategyTitle, SyncPrototypeStrategyHelp, inputs)
	if err != nil || aborted {
		return configdomain.SyncPrototypeStrategyMerge, aborted, err
	}
	cutSelection, _, _ := strings.Cut(selection.String(), " ")
	fmt.Printf(messages.SyncPrototypeBranches, components.FormattedSelection(cutSelection, aborted))
	return selection.SyncPrototypeStrategy(), aborted, err
}

type syncPrototypeStrategyEntry string

func (self syncPrototypeStrategyEntry) String() string {
	return string(self)
}

func (self syncPrototypeStrategyEntry) SyncPrototypeStrategy() configdomain.SyncPrototypeStrategy {
	switch self {
	case syncPrototypeStrategyEntryMerge:
		return configdomain.SyncPrototypeStrategyMerge
	case syncPrototypeStrategyEntryRebase:
		return configdomain.SyncPrototypeStrategyRebase
	case syncPrototypeStrategyEntryCompress:
		return configdomain.SyncPrototypeStrategyCompress
	}
	panic("unhandled syncPrototypeStrategyEntry: " + self)
}
