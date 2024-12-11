package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/messages"
)

const (
	newBranchTypeTitle = `New branch type`
	NewBranchTypeHelp  = `
The "new-branch-type" setting determines which branch type Git Town
creates when you run "git town hack", "append", or "prepend".

More info at https://www.git-town.com/preferences/new-branch-type.

`
)

func NewBranchType(existing configdomain.BranchType, inputs components.TestInput) (configdomain.BranchType, bool, error) {
	entries := []configdomain.BranchType{
		configdomain.BranchTypeFeatureBranch,
		configdomain.BranchTypePrototypeBranch,
	}
	var defaultPos int
	if existing == configdomain.BranchTypeFeatureBranch {
		defaultPos = 0
	} else {
		defaultPos = 1
	}
	selection, aborted, err := components.RadioList(list.NewEntries(entries...), defaultPos, newBranchTypeTitle, NewBranchTypeHelp, inputs)
	if err != nil || aborted {
		return configdomain.BranchTypeFeatureBranch, aborted, err
	}
	fmt.Println(messages.CreatePrototypeBranches, components.FormattedSelection(selection.String(), aborted))
	return selection, aborted, err
}
