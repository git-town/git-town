package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	syncPrototypeStrategyTitle = `Sync-prototype strategy`
	SyncPrototypeStrategyHelp  = `
Choose how Git Town should
synchronize prototype branches.

Prototype branches are local-only feature branches.
They are useful for reducing load on CI systems
and limiting the sharing of confidential changes.

`
)

func SyncPrototypeStrategy(args SyncPrototypeStrategyArgs) (Option[configdomain.SyncPrototypeStrategy], dialogdomain.Exit, error) {
	entries := list.Entries[Option[configdomain.SyncPrototypeStrategy]]{
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
	}
	selection, exit, err := ConfigEnumDialog(ConfigEnumDialogArgs[configdomain.SyncPrototypeStrategy]{
		ConfigFileValue: args.ConfigFileValue,
		Entries:         entries,
		HelpText:        SyncPrototypeStrategyHelp,
		Inputs:          args.Inputs,
		LocalValue:      Option[configdomain.SyncPrototypeStrategy]{},
		ParseFunc:       configdomain.ParseSyncPrototypeStrategy,
		Prompt:          "Your sync prototype strategy: ",
		ResultMessage:   messages.SyncPrototypeBranches,
		Title:           syncPrototypeStrategyTitle,
		UnscopedValue:   args.UnscopedValue,
	})
	// selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, syncPrototypeStrategyTitle, SyncPrototypeStrategyHelp, inputs)
	fmt.Printf(messages.SyncPrototypeBranches, dialogcomponents.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}

type SyncPrototypeStrategyArgs struct {
	ConfigFileValue Option[configdomain.SyncPrototypeStrategy]
	Inputs          dialogcomponents.TestInputs
	UnscopedValue   Option[configdomain.SyncPrototypeStrategy]
}
