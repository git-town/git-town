package dialog

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const enterSyncUpstreamHelp = `
Should "git sync" also fetch updates from the upstream remote?

If an "upstream" remote exists, and this setting is enabled,
"git sync" will also update the local main branch
with commits from the main branch at the upstream remote.

This is useful if the repository you work on is a fork,
and you want to keep it in sync with the repo it was forked from.

`

func EnterSyncUpstream(existing configdomain.SyncUpstream, inputs TestInput) (configdomain.SyncUpstream, bool, error) {
	entries := []string{`yes, receive updates from the upstream repo`, `no, don't receive updates from upstream`}
	var defaultPos int
	if existing {
		defaultPos = 0
	} else {
		defaultPos = 1
	}
	selection, aborted, err := radioList(radioListArgs{
		entries:      entries,
		defaultEntry: entries[defaultPos],
		help:         enterSyncUpstreamHelp,
		testInput:    inputs,
	})
	if err != nil || aborted {
		return true, aborted, err
	}
	cutSelection, _, _ := strings.Cut(selection, ",")
	fmt.Printf("Sync with upstream: %s\n", formattedSelection(cutSelection, aborted))
	parsedAnswer, err := configdomain.ParseSyncUpstream(cutSelection, "user dialog")
	return parsedAnswer, aborted, err
}
