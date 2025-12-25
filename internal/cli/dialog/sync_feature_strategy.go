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

func SyncFeatureStrategy(args Args[configdomain.SyncFeatureStrategy]) (Option[configdomain.SyncFeatureStrategy], dialogdomain.Exit, error) {
	entries := list.Entries[Option[configdomain.SyncFeatureStrategy]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[configdomain.SyncFeatureStrategy]]{
			Data: None[configdomain.SyncFeatureStrategy](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[configdomain.SyncFeatureStrategy]]{
		{
			Data: Some(configdomain.SyncFeatureStrategyMerge),
			Text: `merge updates from the parent and tracking branch`,
		},
		{
			Data: Some(configdomain.SyncFeatureStrategyRebase),
			Text: `rebase branches against their parent and tracking branch`,
		},
		{
			Data: Some(configdomain.SyncFeatureStrategyCompress),
			Text: `compress the branch after merging parent and tracking`,
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, syncFeatureStrategyTitle, SyncFeatureStrategyHelp, args.Inputs, "sync-feature-strategy")
	fmt.Printf(messages.SyncFeatureBranches, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
