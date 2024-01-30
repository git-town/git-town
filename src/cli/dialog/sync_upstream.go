package dialog

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/cli/dialog/components"
	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const (
	syncUpstreamTitle = `Sync-upstream strategy`
	syncUpstreamHelp  = `
Should "git sync" also fetch updates from the upstream remote?

If an "upstream" remote exists, and this setting is enabled,
"git sync" will also update the local main branch
with commits from the main branch at the upstream remote.

This is useful if the repository you work on is a fork,
and you want to keep it in sync with the repo it was forked from.

`
)

const (
	SyncUpstreamEntryYes syncUpstreamEntry = `yes, receive updates from the upstream repo`
	SyncUpstreamEntryNo  syncUpstreamEntry = `no, don't receive updates from upstream`
)

func SyncUpstream(existing configdomain.SyncUpstream, inputs components.TestInput) (configdomain.SyncUpstream, bool, error) {
	entries := []syncUpstreamEntry{
		SyncUpstreamEntryYes,
		SyncUpstreamEntryNo,
	}
	var defaultPos int
	if existing {
		defaultPos = 0
	} else {
		defaultPos = 1
	}
	selection, aborted, err := components.RadioList(entries, defaultPos, syncUpstreamTitle, syncUpstreamHelp, inputs)
	if err != nil || aborted {
		return true, aborted, err
	}
	fmt.Printf("Sync with upstream: %s\n", components.FormattedSelection(selection.Short(), aborted))
	return selection.SyncUpstream(), aborted, err
}

type syncUpstreamEntry string

func (self syncUpstreamEntry) Short() string {
	start, _, _ := strings.Cut(self.String(), ",")
	return start
}

func (self syncUpstreamEntry) String() string {
	return string(self)
}

func (self syncUpstreamEntry) SyncUpstream() configdomain.SyncUpstream {
	switch self {
	case SyncUpstreamEntryYes:
		return configdomain.SyncUpstream(true)
	case SyncUpstreamEntryNo:
		return configdomain.SyncUpstream(false)
	}
	panic("unhandled syncUpstreamEntry: " + self)
}
