package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/gohacks/stringers"
)

var PerennialBranchOption = gitdomain.LocalBranchName("<none> (perennial branch)") //nolint:gochecknoglobals

const enterParentHelpTemplate = `
Please select the parent of branch %q or enter its number.
Most of the time this is the main development branch (%v).


`

// EnterParent lets the user select the parent branch for the given branch.
func EnterParent(args EnterParentArgs) (gitdomain.LocalBranchName, bool, error) {
	entries := EnterParentEntries(args)
	cursor := stringers.IndexOrStart(entries, args.MainBranch)
	help := fmt.Sprintf(enterParentHelpTemplate, args.Branch, args.MainBranch)
	selection, aborted, err := radioList(entries, cursor, help, args.DialogTestInput)
	fmt.Printf("Selected parent branch for %q: %s\n", args.Branch, formattedSelection(selection.String(), aborted))
	return selection, aborted, err
}

type EnterParentArgs struct {
	Branch          gitdomain.LocalBranchName
	DialogTestInput TestInput
	LocalBranches   gitdomain.LocalBranchNames
	Lineage         configdomain.Lineage
	MainBranch      gitdomain.LocalBranchName
}

func EnterParentEntries(args EnterParentArgs) gitdomain.LocalBranchNames {
	parentCandidateBranches := args.LocalBranches.Remove(args.Branch).Remove(args.Lineage.Children(args.Branch)...)
	parentCandidateBranches.Sort()
	parentCandidates := parentCandidateBranches.Hoist(args.MainBranch)
	return append(gitdomain.LocalBranchNames{PerennialBranchOption}, parentCandidates...)
}
