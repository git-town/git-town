package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const enterPushHookHelp = `
The "push-hook" setting determines whether Git Town
permits or prevents Git hooks while pushing branches.
Hooks are enabled by default. If your Git hooks are slow,
you can disable them to speed up branch syncing.

When disabled, Git Town pushes using the "--no-verify" switch.
More info at https://www.git-town.com/preferences/push-hook.

`

// EnterMainBranch lets the user select a new main branch for this repo.
func EnterPushHook(existing configdomain.PushHook, inputs TestInput) (configdomain.PushHook, bool, error) {
	var defaultEntry string
	if existing {
		defaultEntry = "yes (default)"
	} else {
		defaultEntry = "no"
	}
	selection, aborted, err := radioList(radioListArgs{
		entries:      []string{"enabled", "disabled"},
		defaultEntry: defaultEntry,
		help:         enterPushHookHelp,
		testInput:    inputs,
	})
	if err != nil || aborted {
		return true, aborted, err
	}
	fmt.Printf("Push hook: %s\n", formattedSelection(selection, aborted))
	result, err := configdomain.NewPushHook(selection, "user dialog")
	return result, aborted, err
}
