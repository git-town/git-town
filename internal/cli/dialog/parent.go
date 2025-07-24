package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents"
	"github.com/git-town/git-town/v21/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/slice"
)

var PerennialBranchOption = gitdomain.LocalBranchName("<none> (perennial branch)")

const (
	parentBranchTitleTemplate = `Parent branch for %s`
	parentBranchHelpTemplate  = `
Please select the parent of branch %q
or enter its number.


`
)

// Parent lets the user select the parent branch for the given branch.
// TODO: delete and replace all usages with dialog.SwitchBranch
func Parent(args ParentArgs) (ParentOutcome, gitdomain.LocalBranchName, error) {
	parentCandidates := ParentCandidateNames(args)
	cursor := slice.Index(parentCandidates, args.DefaultChoice).GetOrElse(0)
	title := fmt.Sprintf(parentBranchTitleTemplate, args.Branch)
	help := fmt.Sprintf(parentBranchHelpTemplate, args.Branch)
	selection, exit, err := dialogcomponents.RadioList(list.NewEntries(parentCandidates...), cursor, title, help, args.Inputs, fmt.Sprintf("parent-branch-for-%q", args.Branch))
	if exit {
		return ParentOutcomeExit, selection, err
	}
	if selection == PerennialBranchOption {
		return ParentOutcomePerennialBranch, selection, err
	}
	return ParentOutcomeSelectedParent, selection, err
}

type ParentArgs struct {
	Branch        gitdomain.LocalBranchName
	DefaultChoice gitdomain.LocalBranchName
	Inputs        dialogcomponents.Inputs
	Lineage       configdomain.Lineage
	LocalBranches gitdomain.LocalBranchNames
	MainBranch    gitdomain.LocalBranchName
}

func ParentCandidateNames(args ParentArgs) gitdomain.LocalBranchNames {
	parentCandidateBranches := args.LocalBranches.Remove(args.Branch).Remove(args.Lineage.Children(args.Branch)...)
	parentCandidateBranches = slice.NaturalSort(parentCandidateBranches)
	parentCandidates := parentCandidateBranches.Hoist(args.MainBranch)
	return append(gitdomain.LocalBranchNames{PerennialBranchOption}, parentCandidates...)
}

// ParentOutcome describes the selection that the user made in the `Parent` dialog.
type ParentOutcome int

const (
	ParentOutcomeExit            ParentOutcome = iota // the user exited the dialog
	ParentOutcomePerennialBranch                      // the user chose the "perennial branch" option
	ParentOutcomeSelectedParent                       // the user selected one of the branches
)
