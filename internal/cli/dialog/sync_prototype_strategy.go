package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/messages"
)

const (
	syncPrototypeStrategyTitle = `Sync-prototype strategy`
	SyncPrototypeStrategyHelp  = `
How should Git Town synchronize prototype branches?
Prototype branches are feature branches that haven't been proposed yet.
Typically they contain  features and bug fixes on them,
hence their name.

`
)

func SyncPrototypeStrategy(existing configdomain.SyncPrototypeStrategy, inputs components.TestInput) (configdomain.SyncPrototypeStrategy, bool, error) {
	entries := list.Entries[configdomain.SyncPrototypeStrategy]{
		{
			Data:     configdomain.SyncPrototypeStrategyMerge,
			Disabled: false,
			Text:     "merge updates from the parent and tracking branch",
		},
		{
			Data:     configdomain.SyncPrototypeStrategyRebase,
			Disabled: false,
			Text:     "rebase branches against their parent and tracking branch",
		},
		{
			Data:     configdomain.SyncPrototypeStrategyCompress,
			Disabled: false,
			Text:     "compress the branch after merging parent and tracking",
		},
	}
	defaultPos := entries.IndexOf(existing)
	selection, aborted, err := components.RadioList(entries, defaultPos, syncPrototypeStrategyTitle, SyncPrototypeStrategyHelp, inputs)
	if err != nil || aborted {
		return configdomain.SyncPrototypeStrategyMerge, aborted, err
	}
	fmt.Printf(messages.SyncPrototypeBranches, components.FormattedSelection(selection.String(), aborted))
	return selection, aborted, err
}
