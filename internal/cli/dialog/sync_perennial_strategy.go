package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
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

func SyncPerennialStrategy(existing configdomain.SyncPerennialStrategy, inputs dialogcomponents.TestInput) (configdomain.SyncPerennialStrategy, dialogdomain.Exit, error) {
	entries := list.Entries[configdomain.SyncPerennialStrategy]{
		{
			Data: configdomain.SyncPerennialStrategyFFOnly,
			Text: "fast-forward perennial branches to their tracking branch",
		},
		{
			Data: configdomain.SyncPerennialStrategyRebase,
			Text: "rebase perennial branches against their tracking branch",
		},
	}
	defaultPos := entries.IndexOf(existing)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, syncPerennialStrategyTitle, SyncPerennialStrategyHelp, inputs)
	fmt.Printf(messages.SyncPerennialBranches, dialogcomponents.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}
