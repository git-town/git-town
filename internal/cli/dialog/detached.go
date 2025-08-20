package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	detachedTitle = `Sync in detached mode`
	detachedHelp  = `
By default, "git town sync" pulls
the latest commits from the main branch
into your feature branches.

Detached mode prevents this.  When enabled,
feature branches sync only with their tracking
and non-perennial parent branches.

To manually pull new commits from the main branch
in detached mode, run "git town sync --no-detached".

`
)

func SyncDetached(args Args[configdomain.Detached]) (Option[configdomain.Detached], dialogdomain.Exit, error) {
	entries := list.Entries[Option[configdomain.Detached]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[configdomain.Detached]]{
			Data: None[configdomain.Detached](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[configdomain.Detached]]{
		{
			Data: Some(configdomain.Detached(false)),
			Text: "yes, pull updates from the main branch when syncing",
		},
		{
			Data: Some(configdomain.Detached(true)),
			Text: "no, I will pull updates from main manually",
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, detachedTitle, detachedHelp, args.Inputs, "detached")
	fmt.Printf(messages.DetachedResult, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
