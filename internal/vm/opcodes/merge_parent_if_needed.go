package opcodes

import (
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/messages"
	"github.com/git-town/git-town/v16/internal/vm/shared"
	. "github.com/git-town/git-town/v16/pkg/prelude"
)

// merges the branch that at runtime is the parent branch of the given branch into the given branch
type MergeParentIfNeeded struct {
	Branch                  gitdomain.LocalBranchName
	OriginalParentName      Option[gitdomain.LocalBranchName]
	OriginalParentSHA       Option[gitdomain.SHA]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *MergeParentIfNeeded) Run(args shared.RunArgs) error {
	program := []shared.Opcode{}
	branchInfos, hasBranchInfos := args.BranchInfos.Get()
	if !hasBranchInfos {
		panic(messages.BranchInfosNotProvided)
	}
	for branch := self.Branch; ; {
		parent, hasParent := args.Config.Value.NormalConfig.Lineage.Parent(branch).Get()
		if !hasParent {
			break
		}
		parentIsLocal := branchInfos.HasLocalBranch(parent)
		if parentIsLocal {
			var parentToMerge gitdomain.BranchName
			if branchInfos.BranchIsActiveInAnotherWorktree(parent) {
				parentToMerge = parent.TrackingBranch().BranchName()
			} else {
				parentToMerge = parent.BranchName()
			}
			program = append(program, &MergeParent{
				CurrentParent:      parentToMerge,
				OriginalParentName: self.OriginalParentName,
				OriginalParentSHA:  self.OriginalParentSHA,
			})
			break
		}
		// here the parent isn't local --> sync with its tracking branch, then try again with the grandparent until we find a local ancestor
		parentIsRemote := branchInfos.HasMatchingTrackingBranchFor(parent)
		if parentIsRemote {
			parentTrackingBranch := parent.AtRemote(gitdomain.RemoteOrigin)
			program = append(program, &MergeParent{
				CurrentParent:      parentTrackingBranch.BranchName(),
				OriginalParentName: self.OriginalParentName,
				OriginalParentSHA:  self.OriginalParentSHA,
			})
		}
		branch = parent
	}
	args.PrependOpcodes(program...)
	return nil
}
