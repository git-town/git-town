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

const (
	pushNewBranchesEntryYes pushNewBranchesEntry = "yes, push new branches to origin"
	pushNewBranchesEntryNo  pushNewBranchesEntry = "no, new branches remain local until synced"
)

func EnterPushNewBranches(existing configdomain.NewBranchPush, inputs TestInput) (configdomain.NewBranchPush, bool, error) {
	entries := []pushNewBranchesEntry{pushNewBranchesEntryYes, pushNewBranchesEntryNo}
	var defaultPos int
	if existing {
		defaultPos = 0
	} else {
		defaultPos = 1
	}
	selection, aborted, err := radioList(entries, defaultPos, enterPushNewBranchesHelp, inputs)
	if err != nil || aborted {
		return true, aborted, err
	}
	fmt.Printf("Push new branches: %s\n", formattedSelection(selection.Short(), aborted))
	return selection.ToNewBranchPush(), aborted, err
}

type pushNewBranchesEntry string

func (self pushNewBranchesEntry) Short() string {
	begin, _, _ := strings.Cut(self.String(), ",")
	return begin
}

func (self pushNewBranchesEntry) String() string {
	return string(self)
}

func (self pushNewBranchesEntry) ToNewBranchPush() configdomain.NewBranchPush {
	switch self {
	case pushNewBranchesEntryYes:
		return configdomain.NewBranchPush(true)
	case pushNewBranchesEntryNo:
		return configdomain.NewBranchPush(false)
	}
	panic("unhandled pushNewBranchesEntry: " + self)
}
