package sync

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/vm/opcodes"
)

// PerennialBranchProgram adds the opcodes to sync the perennial branch with the given name.
func PerennialBranchProgram(branch gitdomain.BranchInfo, args BranchProgramArgs) {
	if remoteBranch, hasRemoteBranch := branch.RemoteName.Get(); hasRemoteBranch {
		updateCurrentPerennialBranchOpcode(args.Program, remoteBranch, args.Config.NormalConfig.SyncPerennialStrategy)
	}
	if localBranch, hasLocalBranch := branch.LocalName.Get(); hasLocalBranch {
		if localBranch == args.Config.ValidatedConfigData.MainBranch && args.Remotes.HasUpstream() && args.Config.NormalConfig.SyncUpstream.IsTrue() {
			args.Program.Value.Add(&opcodes.FetchUpstream{Branch: args.Config.ValidatedConfigData.MainBranch})
			args.Program.Value.Add(&opcodes.RebaseBranch{Branch: gitdomain.NewBranchName("upstream/" + args.Config.ValidatedConfigData.MainBranch.String())})
		}
	}
}
