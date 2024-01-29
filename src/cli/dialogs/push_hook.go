package dialogs

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/dialogs/components"
	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const pushHookHelp = `
The "push-hook" setting determines whether Git Town
permits or prevents Git hooks while pushing branches.
Hooks are enabled by default. If your Git hooks are slow,
you can disable them to speed up branch syncing.

When disabled, Git Town pushes using the "--no-verify" switch.
More info at https://www.git-town.com/preferences/push-hook.

Push hooks in Git Town are:
`

const (
	pushHookEntryEnabled  pushHookEntry = "enabled"
	pushHookEntryDisabled pushHookEntry = "disabled"
)

func PushHook(existing configdomain.PushHook, inputs components.TestInput) (configdomain.PushHook, bool, error) {
	entries := []pushHookEntry{
		pushHookEntryEnabled,
		pushHookEntryDisabled,
	}
	var defaultPos int
	if existing {
		defaultPos = 0
	} else {
		defaultPos = 1
	}
	selection, aborted, err := components.RadioList(entries, defaultPos, pushHookHelp, inputs)
	if err != nil || aborted {
		return true, aborted, err
	}
	fmt.Printf("Push hook: %s\n", components.FormattedSelection(selection.String(), aborted))
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
