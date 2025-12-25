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
	syncPerennialStrategyTitle = `Perennial branch sync strategy`
	SyncPerennialStrategyHelp  = `
Choose how Git Town should
synchronize perennial branches.

These branches have no parent
and are only updated
by shipping feature branches.

`
)

func SyncPerennialStrategy(args Args[configdomain.SyncPerennialStrategy]) (Option[configdomain.SyncPerennialStrategy], dialogdomain.Exit, error) {
	entries := list.Entries[Option[configdomain.SyncPerennialStrategy]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[configdomain.SyncPerennialStrategy]]{
			Data: None[configdomain.SyncPerennialStrategy](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[configdomain.SyncPerennialStrategy]]{
		{
			Data: Some(configdomain.SyncPerennialStrategyFFOnly),
			Text: "fast-forward perennial branches to their tracking branch",
		},
		{
			Data: Some(configdomain.SyncPerennialStrategyRebase),
			Text: "rebase perennial branches against their tracking branch",
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, syncPerennialStrategyTitle, SyncPerennialStrategyHelp, args.Inputs, "sync-perennial-strategy")
	fmt.Printf(messages.SyncPerennialBranches, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
