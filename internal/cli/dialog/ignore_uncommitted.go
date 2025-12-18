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
	ignoreUncommittedTitle = `Ignore uncommitted changes when shipping`
	IgnoreUncommittedHelp  = `
Should git town ship proceed
when there are uncommitted changes?

By default, shipping branches with uncommitted changes will final
to ensure all changes on the branch are shipped.

`
)

func IgnoreUncommitted(args Args[configdomain.IgnoreUncommitted]) (Option[configdomain.IgnoreUncommitted], dialogdomain.Exit, error) {
	entries := list.Entries[Option[configdomain.IgnoreUncommitted]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[configdomain.IgnoreUncommitted]]{
			Data: None[configdomain.IgnoreUncommitted](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[configdomain.IgnoreUncommitted]]{
		{
			Data: Some(configdomain.IgnoreUncommitted(true)),
			Text: `yes, ship branches with uncommitted changes`,
		},
		{
			Data: Some(configdomain.IgnoreUncommitted(false)),
			Text: `no, remind me to commit changes`,
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, ignoreUncommittedTitle, IgnoreUncommittedHelp, args.Inputs, "ignore-uncommitted")
	fmt.Printf(messages.IgnoreUncommitted, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
