package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v17/internal/cli/dialog/components"
	"github.com/git-town/git-town/v17/internal/cli/dialog/components/list"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/gohacks/slice"
	"github.com/git-town/git-town/v17/internal/messages"
)

var PerennialBranchOption = gitdomain.LocalBranchName("<none> (perennial branch)") //nolint:gochecknoglobals

const (
	parentBranchTitleTemplate = `Parent branch for %s`
	parentBranchHelpTemplate  = `
Please select the parent of branch %q or enter its number.
Most of the time this is the main branch (%v).


`
)

// Parent lets the user select the parent branch for the given branch.
func Parent(args ParentArgs) (ParentOutcome, gitdomain.LocalBranchName, error) {
	parentCandidates := ParentCandidateNames(args)
	cursor := slice.Index(parentCandidates, args.DefaultChoice).GetOrElse(0)
	title := fmt.Sprintf(parentBranchTitleTemplate, args.Branch)
	help := fmt.Sprintf(parentBranchHelpTemplate, args.Branch, args.MainBranch)
	selection, aborted, err := components.RadioList(list.NewEntries(parentCandidates...), cursor, title, help, args.DialogTestInput)
	fmt.Printf(messages.ParentDialogSelected, args.Branch, components.FormattedSelection(selection.String(), aborted))
	if aborted {
		return ParentOutcomeAborted, selection, err
	}
	if selection == PerennialBranchOption {
		return ParentOutcomePerennialBranch, selection, err
	}
	return ParentOutcomeSelectedParent, selection, err
}

type ParentArgs struct {
	Branch          gitdomain.LocalBranchName
	DefaultChoice   gitdomain.LocalBranchName
	DialogTestInput components.TestInput
	Lineage         configdomain.Lineage
	LocalBranches   gitdomain.LocalBranchNames
	MainBranch      gitdomain.LocalBranchName
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
	ParentOutcomeAborted         ParentOutcome = iota // the user aborted the dialog
	ParentOutcomePerennialBranch                      // the user chose the "perennial branch" option
	ParentOutcomeSelectedParent                       // the user selected one of the branches
)
