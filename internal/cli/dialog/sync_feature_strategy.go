package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/messages"
)

const (
	syncFeatureStrategyTitle = `Sync-feature strategy`
	SyncFeatureStrategyHelp  = `
How should Git Town synchronize feature branches?
Feature branches are short-lived branches cut from
the main branch and shipped back into the main branch.
Typically you develop features and bug fixes on them,
hence their name.

`
)

func SyncFeatureStrategy(existing configdomain.SyncFeatureStrategy, inputs components.TestInput) (configdomain.SyncFeatureStrategy, bool, error) {
	entries := list.Entries[configdomain.SyncFeatureStrategy]{
		{
			Data:    configdomain.SyncFeatureStrategyMerge,
			Enabled: true,
			Text:    `merge updates from the parent and tracking branch`,
		},
		{
			Data:    configdomain.SyncFeatureStrategyRebase,
			Enabled: true,
			Text:    `rebase branches against their parent and tracking branch`,
		},
		{
			Data:    configdomain.SyncFeatureStrategyCompress,
			Enabled: true,
			Text:    `compress the branch after merging parent and tracking`,
		},
	}
	defaultPos := DialogPosition(entries, existing)
	selection, aborted, err := components.RadioList(entries, defaultPos, syncFeatureStrategyTitle, SyncFeatureStrategyHelp, inputs)
	if err != nil || aborted {
		return configdomain.SyncFeatureStrategyMerge, aborted, err
	}
	fmt.Printf(messages.SyncFeatureBranches, components.FormattedSelection(selection.String(), aborted))
	return selection, aborted, err
}
