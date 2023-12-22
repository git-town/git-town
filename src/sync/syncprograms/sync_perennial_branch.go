package syncprograms

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/sync/syncdomain"
	"github.com/git-town/git-town/v11/src/vm/opcode"
)

// SyncPerennialBranchProgram adds the opcodes to sync the perennial branch with the given name.
func SyncPerennialBranchProgram(branch syncdomain.BranchInfo, args SyncBranchProgramArgs) {
	if branch.HasTrackingBranch() {
		updateCurrentPerennialBranchOpcode(args.Program, branch.RemoteName, args.SyncPerennialStrategy)
	}
	if branch.LocalName == args.MainBranch && args.Remotes.HasUpstream() && args.SyncUpstream.Bool() {
		args.Program.Add(&opcode.FetchUpstream{Branch: args.MainBranch})
		args.Program.Add(&opcode.RebaseBranch{Branch: gitdomain.NewBranchName("upstream/" + args.MainBranch.String())})
	}
}
