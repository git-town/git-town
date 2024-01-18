package dialog

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const enterPushNewBranchesHelp = `
Should Git Town push the new branches it creates
immediately to origin even if they are empty?

When enabled, you can run "git push" right away
but creating new branches is slower and
it triggers an unnecessary CI run on the empty branch.

When disabled, many Git Town commands execute faster
and Git Town will create the missing tracking branch
on the first run of "git sync".

`

func EnterPushNewBranches(existing configdomain.NewBranchPush, inputs TestInput) (configdomain.NewBranchPush, bool, error) {
	entries := []string{"yes, push new branches to origin", "no, new branches remain local until synced"}
	var defaultPos int
	if existing {
		defaultPos = 0
	} else {
		defaultPos = 1
	}
	selection, aborted, err := radioList(radioListArgs{
		entries:      entries,
		defaultEntry: entries[defaultPos],
		help:         enterPushNewBranchesHelp,
		testInput:    inputs,
	})
	if err != nil || aborted {
		return true, aborted, err
	}
	fmt.Printf("Push new branches: %s\n", formattedSelection(selection, aborted))
	cutSelection, _, _ := strings.Cut(selection, ",")
	parsedAnswer, err := configdomain.ParseNewBranchPush(cutSelection, "user dialog")
	return parsedAnswer, aborted, err
}
