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
	interactiveTitle = `Interactive mode`
	interactiveHelp  = `
Git Town prompts for missing input if needed.

These features require an interactive terminal.
This setting allows to disable interactive mode.

More details: https://www.git-town.com/preferences/interactive.

`
)

func Interactive(args Args[configdomain.Interactive]) (Option[configdomain.Interactive], dialogdomain.Exit, error) {
	entries := list.Entries[Option[configdomain.Interactive]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[configdomain.Interactive]]{
			Data: None[configdomain.Interactive](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[configdomain.Interactive]]{
		{
			Data: Some(configdomain.InteractiveEnabled),
			Text: "enabled: prompt for missing input",
		},
		{
			Data: Some(configdomain.Interactive("disabled")),
			Text: "disabled: enter all input via CLI flags",
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, interactiveTitle, interactiveHelp, args.Inputs, args.Interactive, "interactive")
	fmt.Printf(messages.Interactive, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
