package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v17/internal/cli/format"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/messages"
)

const (
	syncUpstreamTitle = `Sync-upstream strategy`
	SyncUpstreamHelp  = `
Should "git town sync" also fetch updates from the upstream remote?

If an "upstream" remote exists, and this setting is enabled,
"git town sync" will also update the local main branch
with commits from the main branch at the upstream remote.

This is useful if the repository you work on is a fork,
and you want to keep it in sync with the repo it was forked from.

`
)

func SyncUpstream(existing configdomain.SyncUpstream, inputs components.TestInput) (configdomain.SyncUpstream, bool, error) {
	entries := list.Entries[configdomain.SyncUpstream]{
		{
			Data:    true,
			Enabled: true,
			Text:    "yes, receive updates from the upstream repo",
		},
		{
			Data:    false,
			Enabled: true,
			Text:    "no, don't receive updates from upstream",
		},
	}
	defaultPos := entries.IndexOfData(existing)
	selection, aborted, err := components.RadioList(entries, defaultPos, syncUpstreamTitle, SyncUpstreamHelp, inputs)
	if err != nil || aborted {
		return true, aborted, err
	}
	fmt.Printf(messages.SyncWithUpstream, components.FormattedSelection(format.Bool(selection.IsTrue()), aborted))
	return selection, aborted, err
}
