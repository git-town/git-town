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
	autoSyncTitle = `Auto-sync`
	autoSyncHelp  = `
Git Town can keep your branches in sync,
for example before creating a new branch.

If you disable this, branches only get synced
if you run "git town sync" manually.
`
)

func AutoSync(args Args[configdomain.AutoSync]) (Option[configdomain.AutoSync], dialogdomain.Exit, error) {
	entries := list.Entries[Option[configdomain.AutoSync]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[configdomain.AutoSync]]{
			Data: None[configdomain.AutoSync](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[configdomain.AutoSync]]{
		{
			Data: Some(configdomain.AutoSync(true)),
			Text: "yes, sync branches if possible",
		},
		{
			Data: Some(configdomain.AutoSync(false)),
			Text: `no, I'll run "git town sync" to sync my branches`,
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, autoSyncTitle, autoSyncHelp, args.Inputs, "auto-sync")
	fmt.Printf(messages.AutoSync, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
