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
	detachedTitle = `Sync detached`
	detachedHelp  = `
Should "git town sync" pull
updates from the main branch
into feature branches?

Disabling this makes sense if
too frequent unrelated updates
in busy monorepos are a problem.

When disabled, you would then need to
pull updates manually by running
git town sync --detached=0

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
