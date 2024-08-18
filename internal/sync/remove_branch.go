package sync

import (
	"github.com/git-town/git-town/v15/internal/config/configdomain"
	"github.com/git-town/git-town/v15/internal/git/gitdomain"
	"github.com/git-town/git-town/v15/internal/vm/opcodes"
	"github.com/git-town/git-town/v15/internal/vm/program"
	. "github.com/git-town/git-town/v15/pkg/prelude"
)

func RemoveBranchFromLineage(args RemoveBranchFromLineageArgs) {
	childBranches := args.Lineage.Children(args.Branch)
	for _, child := range childBranches {
		args.Program.Value.Add(&opcodes.ChangeParent{Branch: child, Parent: args.Parent})
	}
	args.Program.Value.Add(&opcodes.DeleteParentBranch{Branch: args.Branch})
}

type RemoveBranchFromLineageArgs struct {
	Branch  gitdomain.LocalBranchName
	Lineage configdomain.Lineage
	Parent  gitdomain.LocalBranchName
	Program Mutable[program.Program]
}
