package sync

import (
	"github.com/git-town/git-town/v13/src/config/configdomain"
	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/src/vm/opcodes"
	"github.com/git-town/git-town/v13/src/vm/program"
)

func RemoveBranchFromLineage(args RemoveBranchFromLineageArgs) {
	childBranches := args.Lineage.Children(args.Branch)
	for _, child := range childBranches {
		args.Program.Add(&opcodes.ChangeParent{Branch: child, Parent: args.Parent})
	}
	args.Program.Add(&opcodes.DeleteParentBranch{Branch: args.Branch})
}

type RemoveBranchFromLineageArgs struct {
	Branch  gitdomain.LocalBranchName
	Lineage configdomain.Lineage
	Parent  gitdomain.LocalBranchName
	Program *program.Program
}
