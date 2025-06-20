package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v21/internal/cli/format"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/messages"
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

func PushHook(existing configdomain.PushHook, inputs components.TestInput) (configdomain.PushHook, dialogdomain.Aborted, error) {
	entries := list.Entries[configdomain.PushHook]{
		{
			Data: true,
			Text: "enabled: run Git hooks when pushing branches",
		},
		{
			Data: false,
			Text: "disabled: don't run Git hooks when pushing branches",
		},
	}
	defaultPos := entries.IndexOf(existing)
	selection, aborted, err := components.RadioList(entries, defaultPos, pushHookTitle, PushHookHelp, inputs)
	if err != nil || aborted {
		return true, aborted, err
	}
	fmt.Printf(messages.PushHook, components.FormattedSelection(format.Bool(selection.IsTrue()), aborted))
	return selection, aborted, err
}
