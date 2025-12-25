package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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

func SyncUpstream(args Args[configdomain.SyncUpstream]) (Option[configdomain.SyncUpstream], dialogdomain.Exit, error) {
	entries := list.Entries[Option[configdomain.SyncUpstream]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[configdomain.SyncUpstream]]{
			Data: None[configdomain.SyncUpstream](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[configdomain.SyncUpstream]]{
		{
			Data: Some(configdomain.SyncUpstream(true)),
			Text: "yes, receive updates from the upstream repo",
		},
		{
			Data: Some(configdomain.SyncUpstream(false)),
			Text: "no, don't receive updates from upstream",
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, syncUpstreamTitle, SyncUpstreamHelp, args.Inputs, "sync-upstream")
	fmt.Printf(messages.SyncWithUpstream, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
