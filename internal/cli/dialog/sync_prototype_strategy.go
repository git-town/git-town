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
	syncPrototypeStrategyTitle = `Prototype branch sync strategy`
	SyncPrototypeStrategyHelp  = `
Choose how Git Town should
synchronize prototype branches.

Prototype branches are local-only feature branches.
They are useful for reducing load on CI systems
and limiting the sharing of confidential changes.

`
)

func SyncPrototypeStrategy(args CommonArgs) (Option[configdomain.SyncPrototypeStrategy], dialogdomain.Exit, error) {
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
		ConfigFileValue: args.ConfigFile.SyncPrototypeStrategy,
		Entries:         entries,
		HelpText:        SyncPrototypeStrategyHelp,
		Inputs:          args.Inputs,
		LocalValue:      args.LocalGitConfig.SyncPrototypeStrategy,
		ParseFunc:       configdomain.ParseSyncPrototypeStrategy,
		Prompt:          "Your sync prototype strategy: ",
		ResultMessage:   messages.SyncPrototypeBranches,
		Title:           syncPrototypeStrategyTitle,
		UnscopedValue:   args.UnscopedGitConfig.SyncPrototypeStrategy,
	})
	// selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, syncPrototypeStrategyTitle, SyncPrototypeStrategyHelp, inputs)
	fmt.Printf(messages.SyncPrototypeBranches, dialogcomponents.FormattedSelection(selection.String(), exit))
	return selection, exit, err
}
