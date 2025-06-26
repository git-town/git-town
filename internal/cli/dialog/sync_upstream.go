package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/cli/format"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
)

const (
	syncUpstreamTitle = `Sync with upstream remote`
	SyncUpstreamHelp  = `
Should git town sync also fetch
updates from the upstream remote?

If an upstream remote exists
and this option is enabled,
"git town sync" will update
your local main branch
with commits from the
upstream's main branch.

This keeps a forked repository
up-to-date with changes
made to the original project.

`
)

func SyncUpstream(existing configdomain.SyncUpstream, inputs components.TestInput) (configdomain.SyncUpstream, dialogdomain.Exit, error) {
	entries := list.Entries[configdomain.SyncUpstream]{
		{
			Data: true,
			Text: "yes, receive updates from the upstream repo",
		},
		{
			Data: false,
			Text: "no, don't receive updates from upstream",
		},
	}
	defaultPos := entries.IndexOf(existing)
	selection, exit, err := components.RadioList(entries, defaultPos, syncUpstreamTitle, SyncUpstreamHelp, inputs)
	if err != nil || exit {
		return true, exit, err
	}
	fmt.Printf(messages.SyncWithUpstream, components.FormattedSelection(format.Bool(selection.IsTrue()), exit))
	return selection, exit, err
}
