package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v19/internal/cli/dialog/components"
	"github.com/git-town/git-town/v19/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v19/internal/config/configdomain"
	"github.com/git-town/git-town/v19/internal/messages"
)

const (
	syncPerennialStrategyTitle = `Perennial branch sync strategy`
	SyncPerennialStrategyHelp  = `
Choose how Git Town should synchronize perennial branches.

These branches have no parent and are only updated
via new commits pushed to their tracking branch from elsewhere.

`
)

func SyncPerennialStrategy(existing configdomain.SyncPerennialStrategy, inputs components.TestInput) (configdomain.SyncPerennialStrategy, bool, error) {
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
	selection, aborted, err := components.RadioList(entries, defaultPos, syncPerennialStrategyTitle, SyncPerennialStrategyHelp, inputs)
	if err != nil || aborted {
		return configdomain.SyncPerennialStrategyRebase, aborted, err
	}
	fmt.Printf(messages.SyncPerennialBranches, components.FormattedSelection(selection.String(), aborted))
	return selection, aborted, err
}
