package dialog

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/cli/dialog/components"
	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const syncFeatureStrategyHelp = `
How should Git Town synchronize feature branches?

Feature branches are short-lived branches cut from the main branch
and shipped back into the main branch.
Typically you develop features and bug fixes on them, hence their name.

How should Git Town update feature branches?

`

const (
	syncFeatureStrategyEntryMerge  syncFeatureStrategyEntry = `merge updates from the parent branch into feature branches`
	syncFeatureStrategyEntryRebase syncFeatureStrategyEntry = `rebase feature branches against their parent branch`
)

func SyncFeatureStrategy(existing configdomain.SyncFeatureStrategy, inputs components.TestInput) (configdomain.SyncFeatureStrategy, bool, error) {
	entries := []syncFeatureStrategyEntry{
		syncFeatureStrategyEntryMerge,
		syncFeatureStrategyEntryRebase,
	}
	var defaultPos int
	switch existing {
	case configdomain.SyncFeatureStrategyMerge:
		defaultPos = 0
	case configdomain.SyncFeatureStrategyRebase:
		defaultPos = 1
	default:
		panic("unknown sync-feature-strategy: " + existing.String())
	}
	selection, aborted, err := components.RadioList(entries, defaultPos, syncFeatureStrategyHelp, inputs)
	if err != nil || aborted {
		return configdomain.SyncFeatureStrategyMerge, aborted, err
	}
	cutSelection, _, _ := strings.Cut(selection.String(), " ")
	fmt.Printf("Sync feature branches: %s\n", components.FormattedSelection(cutSelection, aborted))
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
	}
	panic("unhandled syncFeatureStrategyEntry: " + self)
}
