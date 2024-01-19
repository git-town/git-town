package dialog

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const enterSyncUpstreamHelp = `
Should "git sync" also fetch updates from the upstream remote?

If an "upstream" remote exists, and this setting is enabled,
"git sync" will also rebase/merge the local main branch
against that branch at the upstream remote.

This is useful if the repository you work on is a fork,
and you want to keep it in sync with the repo it was forked from.

`

func enterSyncUpstream(existing configdomain.SyncUpstream, inputs TestInput) (configdomain.SyncUpstream, bool, error) {
	entries := []string{`yes, "git sync" should also pull`, `no, "git ship" should not sync the branch`}
	var defaultPos int
	if existing {
		defaultPos = 0
	} else {
		defaultPos = 1
	}
	selection, aborted, err := radioList(radioListArgs{
		entries:      entries,
		defaultEntry: entries[defaultPos],
		help:         enterSyncBeforeShipHelp,
		testInput:    inputs,
	})
	if err != nil || aborted {
		return true, aborted, err
	}
	cutSelection, _, _ := strings.Cut(selection, ",")
	fmt.Printf("Sync before ship: %s\n", formattedSelection(cutSelection, aborted))
	parsedAnswer, err := configdomain.ParseSyncBeforeShip(cutSelection, "user dialog")
	return parsedAnswer, aborted, err
}
