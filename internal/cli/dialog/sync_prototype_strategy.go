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
	syncPrototypeStrategyTitle = `Prototype branch sync strategy`
	SyncPrototypeStrategyHelp  = `
Choose how Git Town should
synchronize prototype branches.

Prototype branches are local-only feature branches.
They are useful for reducing load on CI systems
and limiting the sharing of confidential changes.

`
)

func SyncPrototypeStrategy(args Args[configdomain.SyncPrototypeStrategy]) (Option[configdomain.SyncPrototypeStrategy], dialogdomain.Exit, error) {
	entries := list.Entries[Option[configdomain.SyncPrototypeStrategy]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[configdomain.SyncPrototypeStrategy]]{
			Data: None[configdomain.SyncPrototypeStrategy](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[configdomain.SyncPrototypeStrategy]]{
		{
			Data: Some(configdomain.SyncPrototypeStrategyMerge),
			Text: "merge updates from the parent and tracking branch",
		},
		{
			Data: Some(configdomain.SyncPrototypeStrategyRebase),
			Text: "rebase branches against their parent and tracking branch",
		},
		{
			Data: Some(configdomain.SyncPrototypeStrategyCompress),
			Text: "compress the branch after merging parent and tracking",
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, syncPrototypeStrategyTitle, SyncPrototypeStrategyHelp, args.Inputs, "sync-prototype-strategy")
	fmt.Printf(messages.SyncPrototypeBranches, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
