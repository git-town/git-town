package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/cli/format"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

const (
	pushHookTitle = `Push hook`
	PushHookHelp  = `
The push-hook setting controls whether
Git Town allows Git hooks to run when pushing branches.
Hooks are enabled by default.
If your Git hooks are slow,
you can disable them to speed up branch syncing.

When disabled, Git Town pushes with the --no-verify flag.

More details: https://www.git-town.com/preferences/push-hook.

`
)

func PushHook(existing Option[configdomain.PushHook], inputs dialogcomponents.TestInputs) (Option[configdomain.PushHook], dialogdomain.Exit, error) {
	entries := list.Entries[Option[configdomain.PushHook]]{
		{
			Data: None[configdomain.PushHook](),
			Text: "not sure, skip this",
		},
		{
			Data: Some(configdomain.PushHook(true)),
			Text: "enabled: run Git hooks when pushing branches",
		},
		{
			Data: Some(configdomain.PushHook(false)),
			Text: "disabled: don't run Git hooks when pushing branches",
		},
	}
	defaultPos := entries.IndexOf(existing)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, pushHookTitle, PushHookHelp, inputs)
	fmt.Printf(messages.PushHook, dialogcomponents.FormattedSelection(format.BoolOpt(selection.Get()), exit))
	return selection, exit, err
}
