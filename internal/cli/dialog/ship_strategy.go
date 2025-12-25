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
	shipStrategyTitle = `Ship strategy`
	ShipStrategyHelp  = `
Choose how Git Town should ship feature branches.

All options update proposals of child branches
and remove the shipped branch locally and remotely.

Options:

- api: merge the proposal on your forge via the forge API
- always-merge: on your machine, merge by
                always creating a merge comment
								(git merge --no-ff)
- fast-forward: on your machine, fast-forward the parent branch
                to include the feature branch commits
- squash-merge: on your machine, squash all commits
                on the feature branch into a single commit
								on the parent branch

`
)

func ShipStrategy(args Args[configdomain.ShipStrategy]) (Option[configdomain.ShipStrategy], dialogdomain.Exit, error) {
	entries := list.Entries[Option[configdomain.ShipStrategy]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[configdomain.ShipStrategy]]{
			Data: None[configdomain.ShipStrategy](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[configdomain.ShipStrategy]]{
		{
			Data: Some(configdomain.ShipStrategyAPI),
			Text: `api: merge the proposal on your forge via the forge API`,
		},
		{
			Data: Some(configdomain.ShipStrategyAlwaysMerge),
			Text: `always-merge: in your local repo, merge the feature branch into its parent by always creating a merge comment (merge --no-ff)`,
		},
		{
			Data: Some(configdomain.ShipStrategyFastForward),
			Text: `fast-forward: in your local repo, fast-forward the parent branch to point to the commits on the feature branch`,
		},
		{
			Data: Some(configdomain.ShipStrategySquashMerge),
			Text: `squash-merge: in your local repo, squash-merge the feature branch into its parent branch`,
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, shipStrategyTitle, ShipStrategyHelp, args.Inputs, "ship-strategy")
	fmt.Printf(messages.ShipStrategy, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
