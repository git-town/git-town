package sync

import (
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	"github.com/git-town/git-town/v24/internal/vm/opcodes"
)

// PerennialBranchProgram adds the opcodes to sync the perennial branch with the given name.
func PerennialBranchProgram(branchInfo gitdomain.BranchInfo, args BranchProgramArgs) {
	if remoteBranch, hasRemoteBranch := branchInfo.RemoteName.Get(); hasRemoteBranch {
		if branchInfo.SyncStatus != gitdomain.SyncStatusUpToDate {
			updateCurrentPerennialBranchOpcode(args.Program, remoteBranch, args.Config.NormalConfig.SyncPerennialStrategy)
		}
	}
	if localBranch, hasLocalBranch := branchInfo.LocalName().Get(); hasLocalBranch {
		if localBranch == args.Config.ValidatedConfigData.MainBranch && args.Remotes.HasUpstream() && args.Config.NormalConfig.SyncUpstream.ShouldSyncUpstream() {
			args.Program.Value.Add(&opcodes.FetchUpstream{Branch: args.Config.ValidatedConfigData.MainBranch})
			args.Program.Value.Add(&opcodes.RebaseBranch{Branch: gitdomain.BranchNameOrPanic(stringss.Trimmed("upstream/" + args.Config.ValidatedConfigData.MainBranch.String()))}) // this string will never need to be trimmed
		}
	}
}
