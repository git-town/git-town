package enter

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/cli/dialogs/components"
	"github.com/git-town/git-town/v11/src/config/configdomain"
)

const pushNewBranchesHelp = `
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
	PushNewBranchesEntryYes pushNewBranchesEntry = "yes, push new branches to origin"
	PushNewBranchesEntryNo  pushNewBranchesEntry = "no, new branches remain local until synced"
)

func PushNewBranches(existing configdomain.NewBranchPush, inputs components.TestInput) (configdomain.NewBranchPush, bool, error) {
	entries := []pushNewBranchesEntry{
		PushNewBranchesEntryYes,
		PushNewBranchesEntryNo,
	}
	var defaultPos int
	if existing {
		defaultPos = 0
	} else {
		defaultPos = 1
	}
	selection, aborted, err := components.RadioList(entries, defaultPos, pushNewBranchesHelp, inputs)
	if err != nil || aborted {
		return true, aborted, err
	}
	fmt.Printf("Push new branches: %s\n", components.FormattedSelection(selection.Short(), aborted))
	return selection.NewBranchPush(), aborted, err
}

type pushNewBranchesEntry string

func (self pushNewBranchesEntry) NewBranchPush() configdomain.NewBranchPush {
	switch self {
	case PushNewBranchesEntryYes:
		return configdomain.NewBranchPush(true)
	case PushNewBranchesEntryNo:
		return configdomain.NewBranchPush(false)
	}
	panic("unhandled pushNewBranchesEntry: " + self)
}

func (self pushNewBranchesEntry) Short() string {
	begin, _, _ := strings.Cut(self.String(), ",")
	return begin
}

func (self pushNewBranchesEntry) String() string {
	return string(self)
}
