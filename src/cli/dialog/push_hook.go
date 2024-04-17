package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/dialog/components/list"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/messages"
)

const (
	pushHookTitle = `Push hook`
	PushHookHelp  = `
The "push-hook" setting determines whether Git Town
permits or prevents Git hooks while pushing branches.
Hooks are enabled by default. If your Git hooks are slow,
you can disable them to speed up branch syncing.

When disabled, Git Town pushes using the "--no-verify" switch.
More info at https://www.git-town.com/preferences/push-hook.

`
)

const (
	pushHookEntryEnabled  pushHookEntry = "enabled"
	pushHookEntryDisabled pushHookEntry = "disabled"
)

func PushHook(existing configdomain.PushHook, inputs components.TestInput) (configdomain.PushHook, bool, error) {
	entries := list.NewEntries(
		pushHookEntryEnabled,
		pushHookEntryDisabled,
	)
	var defaultPos int
	if existing {
		defaultPos = 0
	} else {
		defaultPos = 1
	}
	selection, aborted, err := components.RadioList(list.NewEntries(entries...), defaultPos, pushHookTitle, PushHookHelp, inputs)
	if err != nil || aborted {
		return true, aborted, err
	}
	fmt.Printf(messages.PushHook, components.FormattedSelection(selection.String(), aborted))
	return selection.PushHook(), aborted, err
}

type pushHookEntry string

func (self pushHookEntry) PushHook() configdomain.PushHook {
	switch self {
	case pushHookEntryEnabled:
		return configdomain.PushHook(true)
	case pushHookEntryDisabled:
		return configdomain.PushHook(false)
	}
	panic("unhandled pushHookEntry: " + self)
}

func (self pushHookEntry) String() string {
	return string(self)
}
