package dialog

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const enterPushNewBranchesHelp = `
Should Git Town push the new branches it creates
immediately to origin even if they are empty?

Doing so makes the full setup available right away.
You can run "git push".
The downside is that the extra network operation
makes certain Git Town commands slower
and triggers an unnecessary CI run.

`

func EnterPushNewBranches(existing configdomain.PushHook, inputs TestInput) (configdomain.PushHook, bool, error) {
	entries := []string{"yes, push new branches to origin", "no, don't push new branches to origin"}
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
	selection, _, _ = strings.Cut(selection, ",")
	fmt.Printf("Push new branches: %s\n", formattedSelection(selection, aborted))
	result, err := configdomain.NewPushHook(selection, "user dialog")
	return result, aborted, err
}
