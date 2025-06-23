package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
)

const (
	syncPerennialStrategyTitle = `Perennial branch sync strategy`
	SyncPerennialStrategyHelp  = `
Choose how Git Town should synchronize perennial branches.

These branches have no parent and are only updated
via new commits pushed to their tracking branch from elsewhere.

`
)

func SyncPerennialStrategy(existing configdomain.SyncPerennialStrategy, inputs components.TestInput) (configdomain.SyncPerennialStrategy, dialogdomain.Exit, error) {
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
	selection, exit, err := components.RadioList(entries, defaultPos, syncPerennialStrategyTitle, SyncPerennialStrategyHelp, inputs)
	if err != nil || exit {
		return configdomain.SyncPerennialStrategyRebase, exit, err
	}
	fmt.Printf(messages.SyncPerennialBranches, components.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}
