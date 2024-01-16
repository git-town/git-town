package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const enterPushHookHelp = `
The "push-hook" setting determines whether Git Town permits or prevents Git hooks
while pushing branches.  If your Git hooks are slow, you might want to disable them
when syncing branches to speed up that process. When disabled, Git Town pushes using
the "--no-verify" option, which disables the "pre-push" hook.
More info at https://git-scm.com/docs/githooks#_pre_push.


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
		entries:      []string{"yes", "no"},
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
