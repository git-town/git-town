package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks/slice"
	"github.com/git-town/git-town/v14/src/gohacks/stringers"
	"github.com/git-town/git-town/v14/src/messages"
)

var PerennialBranchOption = gitdomain.LocalBranchName("<none> (perennial branch)") //nolint:gochecknoglobals

const (
	parentBranchTitleTemplate = `Parent branch for %s`
	parentBranchHelpTemplate  = `
Please select the parent of branch %q or enter its number.
Most of the time this is the main development branch (%v).


`
)

// Parent lets the user select the parent branch for the given branch.
func Parent(args ParentArgs) (gitdomain.LocalBranchName, bool, error) {
	entries := ParentEntries(args)
	cursor := stringers.IndexOrStart(entries, args.MainBranch)
	title := fmt.Sprintf(parentBranchTitleTemplate, args.Branch)
	help := fmt.Sprintf(parentBranchHelpTemplate, args.Branch, args.MainBranch)
	selection, aborted, err := components.RadioList(entries, cursor, title, help, args.DialogTestInput)
	fmt.Printf(messages.ParentDialogSelected, args.Branch, components.FormattedSelection(selection.String(), aborted))
	return selection, aborted, err
}

type ParentArgs struct {
	Branch          gitdomain.LocalBranchName
	DialogTestInput components.TestInput
	Lineage         configdomain.Lineage
	LocalBranches   gitdomain.LocalBranchNames
	MainBranch      gitdomain.LocalBranchName
}

func ParentEntries(args ParentArgs) gitdomain.LocalBranchNames {
	parentCandidateBranches := args.LocalBranches.Remove(args.Branch).Remove(args.Lineage.Children(args.Branch)...)
	parentCandidateBranches = slice.NaturalSort(parentCandidateBranches)
	parentCandidates := parentCandidateBranches.Hoist(args.MainBranch)
	return append(gitdomain.LocalBranchNames{PerennialBranchOption}, parentCandidates...)
}
