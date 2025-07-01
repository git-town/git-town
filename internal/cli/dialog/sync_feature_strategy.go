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
	syncFeatureStrategyTitle = `Feature branch sync strategy`
	SyncFeatureStrategyHelp  = `
Choose how Git Town should
synchronize feature branches.

These are short-lived branches
created from the main branch
and eventually merged back into it.
Commonly used for developing
new features and bug fixes.

`
)

func SyncFeatureStrategy(existing configdomain.SyncFeatureStrategy, inputs components.TestInput) (configdomain.SyncFeatureStrategy, dialogdomain.Exit, error) {
	entries := list.Entries[configdomain.SyncFeatureStrategy]{
		{
			Data: configdomain.SyncFeatureStrategyMerge,
			Text: `merge updates from the parent and tracking branch`,
		},
		{
			Data: configdomain.SyncFeatureStrategyRebase,
			Text: `rebase branches against their parent and tracking branch`,
		},
		{
			Data: configdomain.SyncFeatureStrategyCompress,
			Text: `compress the branch after merging parent and tracking`,
		},
	}
	defaultPos := entries.IndexOf(existing)
	selection, exit, err := components.RadioList(entries, defaultPos, syncFeatureStrategyTitle, SyncFeatureStrategyHelp, inputs)
	fmt.Printf(messages.SyncFeatureBranches, components.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}
