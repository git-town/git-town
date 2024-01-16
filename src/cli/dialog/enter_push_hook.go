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

Push hooks in Git Town are:
`

// EnterMainBranch lets the user select a new main branch for this repo.
func EnterPushHook(existing configdomain.PushHook, inputs TestInput) (configdomain.PushHook, bool, error) {
	entries := []string{"enabled", "disabled"}
	var defaultPos int
	if existing {
		defaultPos = 0
	} else {
		defaultPos = 1
	}
	selection, aborted, err := radioList(radioListArgs{
		entries:      entries,
		defaultEntry: entries[defaultPos],
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
