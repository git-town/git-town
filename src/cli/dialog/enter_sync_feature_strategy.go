package dialog

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const enterSyncFeatureStrategyHelp = `
How should Git Town synchronize feature branches?

Feature branches are short-lived branches cut from the main branch
and shipped back into the main branch.
Typically you develop features and bug fixes on them, hence their name.

How should Git Town update feature branches?

`

func EnterSyncFeatureStrategy(existing configdomain.SyncFeatureStrategy, inputs TestInput) (configdomain.SyncFeatureStrategy, bool, error) {
	entries := []string{`merge updates from the parent branch into feature branches`, `rebase feature branches against their parent branch`}
	var defaultPos int
	switch existing {
	case configdomain.SyncFeatureStrategyMerge:
		defaultPos = 0
	case configdomain.SyncFeatureStrategyRebase:
		defaultPos = 1
	default:
		panic("unknown sync-feature-strategy: " + existing.String())
	}
	selection, aborted, err := radioList(radioListArgs{
		entries:      entries,
		defaultEntry: entries[defaultPos],
		help:         enterSyncFeatureStrategyHelp,
		testInput:    inputs,
	})
	if err != nil || aborted {
		return configdomain.SyncFeatureStrategyMerge, aborted, err
	}
	cutSelection, _, _ := strings.Cut(selection, " ")
	fmt.Printf("Sync feature branches: %s\n", formattedSelection(cutSelection, aborted))
	parsedAnswer, err := configdomain.NewSyncFeatureStrategy(cutSelection)
	return parsedAnswer, aborted, err
}
