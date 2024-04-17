package dialog

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/messages"
)

const (
	pushNewBranchesTitle = `Push new branches`
	PushNewBranchesHelp  = `
Should Git Town push the new branches it creates
immediately to origin even if they are empty?

When enabled, you can run "git push" right away
but creating new branches is slower and
it triggers an unnecessary CI run on the empty branch.

When disabled, many Git Town commands execute faster
and Git Town will create the missing tracking branch
on the first run of "git sync".

`
)

const (
	PushNewBranchesEntryYes pushNewBranchesEntry = "yes, push new branches to origin"
	PushNewBranchesEntryNo  pushNewBranchesEntry = "no, new branches remain local until synced"
)

func PushNewBranches(existing configdomain.PushNewBranches, inputs components.TestInput) (configdomain.PushNewBranches, bool, error) {
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
	selection, aborted, err := components.RadioList(components.NewEnabledBubbleListEntries(entries), defaultPos, pushNewBranchesTitle, PushNewBranchesHelp, inputs)
	if err != nil || aborted {
		return true, aborted, err
	}
	fmt.Printf(messages.PushNewBranches, components.FormattedSelection(selection.Short(), aborted))
	return selection.PushNewBranches(), aborted, err
}

type pushNewBranchesEntry string

func (self pushNewBranchesEntry) PushNewBranches() configdomain.PushNewBranches {
	switch self {
	case PushNewBranchesEntryYes:
		return configdomain.PushNewBranches(true)
	case PushNewBranchesEntryNo:
		return configdomain.PushNewBranches(false)
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
