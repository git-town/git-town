package sync

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/vm/opcode"
)

// PerennialBranchProgram adds the opcodes to sync the perennial branch with the given name.
func PerennialBranchProgram(branch gitdomain.BranchInfo, args BranchProgramArgs) {
	if branch.HasTrackingBranch() {
		updateCurrentPerennialBranchOpcode(args.Program, branch.RemoteName, args.Config.SyncPerennialStrategy)
	}
	if branch.LocalName == args.Config.MainBranch && args.Remotes.HasUpstream() && args.Config.SyncUpstream.Bool() {
		args.Program.Add(&opcode.FetchUpstream{Branch: args.Config.MainBranch})
		args.Program.Add(&opcode.RebaseBranch{Branch: gitdomain.NewBranchName("upstream/" + args.Config.MainBranch.String())})
	}
}
