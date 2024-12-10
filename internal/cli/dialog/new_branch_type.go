package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/messages"
)

const (
	createPrototypeBranchesTitle = `Create prototype branches`
	CreatePrototypeBranchesHelp  = `
The "create-prototype-branches" setting determines whether Git Town
always creates prototype branches.
Prototype branches sync only locally and don't create a tracking branch
until they are proposed.

More info at https://www.git-town.com/preferences/create-prototype-branches.

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
	selection, aborted, err := components.RadioList(list.NewEntries(entries...), defaultPos, createPrototypeBranchesTitle, CreatePrototypeBranchesHelp, inputs)
	if err != nil || aborted {
		return configdomain.BranchTypeFeatureBranch, aborted, err
	}
	fmt.Println(messages.CreatePrototypeBranches, components.FormattedSelection(selection.String(), aborted))
	return selection, aborted, err
}
