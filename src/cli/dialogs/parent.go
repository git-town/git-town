package dialogs

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/cli/dialogs/components"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/gohacks/stringers"
)

var PerennialBranchOption = gitdomain.LocalBranchName("<none> (perennial branch)") //nolint:gochecknoglobals

const enterParentHelpTemplate = `
Please select the parent of branch %q or enter its number.
Most of the time this is the main development branch (%v).


`

// Parent lets the user select the parent branch for the given branch.
func Parent(args ParentArgs) (gitdomain.LocalBranchName, bool, error) {
	entries := ParentEntries(args)
	cursor := stringers.IndexOrStart(entries, args.MainBranch)
	help := fmt.Sprintf(enterParentHelpTemplate, args.Branch, args.MainBranch)
	selection, aborted, err := components.RadioList(entries, cursor, help, args.DialogTestInput)
	fmt.Printf("Selected parent branch for %q: %s\n", args.Branch, components.FormattedSelection(selection.String(), aborted))
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
	parentCandidateBranches.Sort()
	parentCandidates := parentCandidateBranches.Hoist(args.MainBranch)
	return append(gitdomain.LocalBranchNames{PerennialBranchOption}, parentCandidates...)
}
