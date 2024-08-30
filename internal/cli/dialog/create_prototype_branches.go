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

const (
	createPrototypeBranchesEntryEnabled  createPrototypeBranchesEntry = "enabled"
	createPrototypeBranchesEntryDisabled createPrototypeBranchesEntry = "disabled"
)

func CreatePrototypeBranches(existing configdomain.CreatePrototypeBranches, inputs components.TestInput) (configdomain.CreatePrototypeBranches, bool, error) {
	entries := []createPrototypeBranchesEntry{
		createPrototypeBranchesEntryEnabled,
		createPrototypeBranchesEntryDisabled,
	}
	var defaultPos int
	if existing {
		defaultPos = 0
	} else {
		defaultPos = 1
	}
	selection, aborted, err := components.RadioList(list.NewEntries(entries...), defaultPos, createPrototypeBranchesTitle, CreatePrototypeBranchesHelp, inputs)
	if err != nil || aborted {
		return false, aborted, err
	}
	fmt.Println(messages.CreatePrototypeBranches, components.FormattedSelection(selection.String(), aborted))
	return selection.CreatePrototypeBranches(), aborted, err
}

type createPrototypeBranchesEntry string

func (self createPrototypeBranchesEntry) CreatePrototypeBranches() configdomain.CreatePrototypeBranches {
	switch self {
	case createPrototypeBranchesEntryEnabled:
		return configdomain.CreatePrototypeBranches(true)
	case createPrototypeBranchesEntryDisabled:
		return configdomain.CreatePrototypeBranches(false)
	}
	panic("unhandled createPrototypeBranchesEntry: " + self)
}

func (self createPrototypeBranchesEntry) String() string {
	return string(self)
}
