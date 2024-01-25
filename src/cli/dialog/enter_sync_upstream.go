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
	entries := []syncUpstreamEntry{syncUpstreamEntryYes, syncUpstreamEntryNo}
	var defaultPos int
	if existing {
		defaultPos = 0
	} else {
		defaultPos = 1
	}
	selection, aborted, err := radioList(entries, defaultPos, enterSyncUpstreamHelp, inputs)
	if err != nil || aborted {
		return true, aborted, err
	}
	cutSelection, _, _ := strings.Cut(selection.String(), ",")
	fmt.Printf("Sync with upstream: %s\n", formattedSelection(cutSelection, aborted))
	return selection.ToSyncUpstream(), aborted, err
}

type syncUpstreamEntry string

func (self syncUpstreamEntry) String() string {
	return string(self)
}

func (self syncUpstreamEntry) ToSyncUpstream() configdomain.SyncUpstream {
	switch self {
	case syncUpstreamEntryYes:
		return configdomain.SyncUpstream(true)
	case syncUpstreamEntryNo:
		return configdomain.SyncUpstream(false)
	}
	panic("unhandled syncUpstreamEntry: " + self)
}

const (
	syncUpstreamEntryYes syncUpstreamEntry = `yes, receive updates from the upstream repo`
	syncUpstreamEntryNo  syncUpstreamEntry = `no, don't receive updates from upstream`
)
