package dialog

import (
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
)

const (
	childBranchTitle = `Child branch`
	childBranchHelp  = `
The current branch has multiple child branches.
Please select which child branch to switch to.

`
)

type ChildBranchArgs struct {
	ChildBranches gitdomain.LocalBranchNames
	Inputs        dialogcomponents.Inputs
}

// ChildBranch lets the user select which child branch to switch to.
func ChildBranch(args ChildBranchArgs) (gitdomain.LocalBranchName, dialogdomain.Exit, error) {
	entries := list.NewEntries(args.ChildBranches...)
	return dialogcomponents.RadioList(entries, 0, childBranchTitle, childBranchHelp, args.Inputs, "child-branch")
}
