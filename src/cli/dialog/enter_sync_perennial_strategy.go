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

func EnterSyncPerennialStrategy(existing configdomain.SyncPerennialStrategy, inputs TestInput) (configdomain.SyncPerennialStrategy, bool, error) {
	entries := []string{`merge updates from the tracking branch into perennial branches`, `rebase perennial branches against their tracking branch`}
	var defaultPos int
	switch existing {
	case configdomain.SyncPerennialStrategyMerge:
		defaultPos = 0
	case configdomain.SyncPerennialStrategyRebase:
		defaultPos = 1
	default:
		panic("unknown sync-perennial-strategy: " + existing.String())
	}
	selection, aborted, err := radioList(radioListArgs{
		entries:      entries,
		defaultEntry: entries[defaultPos],
		help:         enterSyncPerennialStrategyHelp,
		testInput:    inputs,
	})
	if err != nil || aborted {
		return configdomain.SyncPerennialStrategyRebase, aborted, err
	}
	cutSelection, _, _ := strings.Cut(selection, " ")
	fmt.Printf("Sync perennial branches: %s\n", formattedSelection(cutSelection, aborted))
	parsedAnswer, err := configdomain.NewSyncPerennialStrategy(cutSelection)
	return parsedAnswer, aborted, err
}
