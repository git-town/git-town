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

func PushHook(args Args[configdomain.PushHook]) (Option[configdomain.PushHook], dialogdomain.Exit, error) {
	entries := list.Entries[Option[configdomain.PushHook]]{}
	if global, hasGlobal := args.Global.Get(); hasGlobal {
		entries = append(entries, list.Entry[Option[configdomain.PushHook]]{
			Data: None[configdomain.PushHook](),
			Text: fmt.Sprintf(messages.DialogUseGlobalValue, global),
		})
	}
	entries = append(entries, list.Entries[Option[configdomain.PushHook]]{
		{
			Data: Some(configdomain.PushHook(true)),
			Text: "enabled: run Git hooks when pushing branches",
		},
		{
			Data: Some(configdomain.PushHook(false)),
			Text: "disabled: don't run Git hooks when pushing branches",
		},
	}...)
	defaultPos := entries.IndexOf(args.Local)
	selection, exit, err := dialogcomponents.RadioList(entries, defaultPos, pushHookTitle, PushHookHelp, args.Inputs, "push-hook")
	fmt.Printf(messages.PushHook, dialogcomponents.FormattedOption(selection, args.Global.IsSome(), exit))
	return selection, exit, err
}
