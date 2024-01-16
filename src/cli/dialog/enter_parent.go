package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

const PerennialBranchOption = "<none> (perennial branch)"

const enterParentHelpTemplate = `
Please select the parent of branch %q or enter its number.
Most of the time this is the main development branch (%v).


`

// EnterParent lets the user select the parent branch for the given branch.
func EnterParent(args EnterParentArgs) (gitdomain.LocalBranchName, bool, error) {
	selection, aborted, err := radioList(radioListArgs{
		entries:      EnterParentEntries(args),
		defaultEntry: args.MainBranch.String(),
		help:         fmt.Sprintf(enterParentHelpTemplate, args.Branch, args.MainBranch),
		testInput:    args.DialogTestInput,
	})
	fmt.Printf("Selected parent branch for %q: %s\n", args.Branch, formattedSelection(selection, aborted))
	return gitdomain.LocalBranchName(selection), aborted, err
}

type EnterParentArgs struct {
	Branch          gitdomain.LocalBranchName
	DialogTestInput TestInput
	LocalBranches   gitdomain.LocalBranchNames
	Lineage         configdomain.Lineage
	MainBranch      gitdomain.LocalBranchName
}

func EnterParentEntries(args EnterParentArgs) []string {
	parentCandidateBranches := args.LocalBranches.Remove(args.Branch).Remove(args.Lineage.Children(args.Branch)...)
	parentCandidateBranches.Sort()
	parentCandidates := parentCandidateBranches.Hoist(args.MainBranch).Strings()
	return append([]string{PerennialBranchOption}, parentCandidates...)
}
