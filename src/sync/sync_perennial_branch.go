package sync

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/vm/opcodes"
)

// PerennialBranchProgram adds the opcodes to sync the perennial branch with the given name.
func PerennialBranchProgram(branch gitdomain.BranchInfo, args BranchProgramArgs) {
	if remoteBranch, hasRemoteBranch := branch.RemoteName.Get(); hasRemoteBranch {
		updateCurrentPerennialBranchOpcode(args.Program, remoteBranch, args.Config.SyncPerennialStrategy)
	}
	if localBranch, hasLocalBranch := branch.LocalName.Get(); hasLocalBranch {
		if localBranch == args.Config.MainBranch && args.Remotes.HasUpstream() && args.Config.SyncUpstream.Bool() {
			args.Program.Value.Add(&opcodes.FetchUpstream{Branch: args.Config.MainBranch})
			args.Program.Value.Add(&opcodes.RebaseBranch{Branch: gitdomain.NewBranchName("upstream/" + args.Config.MainBranch.String())})
		}
	}
}
