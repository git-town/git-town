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
	syncFeatureStrategyTitle = `Sync-feature strategy`
	SyncFeatureStrategyHelp  = `
How should Git Town synchronize feature branches?
Feature branches are short-lived branches cut from
the main branch and shipped back into the main branch.
Typically you develop features and bug fixes on them,
hence their name.

`
)

const (
	syncFeatureStrategyEntryMerge    syncFeatureStrategyEntry = `merge updates from the parent and tracking branch`
	syncFeatureStrategyEntryRebase   syncFeatureStrategyEntry = `rebase branches against their parent and tracking branch`
	syncFeatureStrategyEntryCompress syncFeatureStrategyEntry = `compress the branch after merging parent and tracking`
)

func SyncFeatureStrategy(existing configdomain.SyncFeatureStrategy, inputs components.TestInput) (configdomain.SyncFeatureStrategy, bool, error) {
	entries := []syncFeatureStrategyEntry{
		syncFeatureStrategyEntryMerge,
		syncFeatureStrategyEntryRebase,
		syncFeatureStrategyEntryCompress,
	}
	var defaultPos int
	switch existing {
	case configdomain.SyncFeatureStrategyMerge:
		defaultPos = 0
	case configdomain.SyncFeatureStrategyRebase:
		defaultPos = 1
	case configdomain.SyncFeatureStrategyCompress:
		defaultPos = 2
	default:
		panic("unknown sync strategy: " + existing.String())
	}
	selection, aborted, err := components.RadioList(list.NewEntries(entries...), defaultPos, syncFeatureStrategyTitle, SyncFeatureStrategyHelp, inputs)
	if err != nil || aborted {
		return configdomain.SyncFeatureStrategyMerge, aborted, err
	}
	cutSelection, _, _ := strings.Cut(selection.String(), " ")
	fmt.Printf(messages.SyncFeatureBranches, components.FormattedSelection(cutSelection, aborted))
	return selection.SyncFeatureStrategy(), aborted, err
}

type syncFeatureStrategyEntry string

func (self syncFeatureStrategyEntry) String() string {
	return string(self)
}

func (self syncFeatureStrategyEntry) SyncFeatureStrategy() configdomain.SyncFeatureStrategy {
	switch self {
	case syncFeatureStrategyEntryMerge:
		return configdomain.SyncFeatureStrategyMerge
	case syncFeatureStrategyEntryRebase:
		return configdomain.SyncFeatureStrategyRebase
	case syncFeatureStrategyEntryCompress:
		return configdomain.SyncFeatureStrategyCompress
	}
	panic("unhandled syncFeatureStrategyEntry: " + self)
}
