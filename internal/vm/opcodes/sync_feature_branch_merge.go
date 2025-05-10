package opcodes

import (
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/messages"
	"github.com/git-town/git-town/v20/internal/vm/shared"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// SyncFeatureBranchMerge merges the parent branches of the given branch until a local parent is found.
type SyncFeatureBranchMerge struct {
	Branch                  gitdomain.LocalBranchName
	InitialParentName       Option[gitdomain.LocalBranchName]
	InitialParentSHA        Option[gitdomain.SHA]
	TrackingBranch          Option[gitdomain.RemoteBranchName]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *SyncFeatureBranchMerge) Run(args shared.RunArgs) error {
	program := []shared.Opcode{}
	branchInfos, hasBranchInfos := args.BranchInfos.Get()
	if !hasBranchInfos {
		panic(messages.BranchInfosNotProvided)
	}
	branch := self.Branch
	for {
		parent, hasParent := args.Config.Value.NormalConfig.Lineage.Parent(branch).Get()
		if !hasParent {
			break
		}
		if parentBranchInfo, hasParentInfo := branchInfos.FindLocalOrRemote(parent, args.Config.Value.NormalConfig.DevRemote).Get(); hasParentInfo {
			parentIsLocal := parentBranchInfo.LocalName.IsSome()
			if parentIsLocal {
				var parentToMerge gitdomain.BranchName
				if branchInfos.BranchIsActiveInAnotherWorktree(parent) {
					parentToMerge = parent.TrackingBranch(args.Config.Value.NormalConfig.DevRemote).BranchName()
				} else {
					parentToMerge = parent.BranchName()
				}
				program = append(program, &MergeParentResolvePhantomConflicts{
					CurrentParent:     parentToMerge,
					InitialParentName: self.InitialParentName,
					InitialParentSHA:  self.InitialParentSHA,
				})
				break
			}
			// here the parent isn't local --> sync with its tracking branch if it exists, then try again with the grandparent until we find a local ancestor
			if parentTrackingBranch, parentHasTrackingBranch := parentBranchInfo.RemoteName.Get(); parentHasTrackingBranch {
				program = append(program, &MergeParentResolvePhantomConflicts{
					CurrentParent:     parentTrackingBranch.BranchName(),
					InitialParentName: self.InitialParentName,
					InitialParentSHA:  self.InitialParentSHA,
				})
			}
		}
		branch = parent
	}
	if trackingBranch, hasTrackingBranch := self.TrackingBranch.Get(); hasTrackingBranch {
		program = append(program, &MergeIntoCurrentBranch{BranchToMerge: trackingBranch.BranchName()})
	}
	args.PrependOpcodes(program...)
	return nil
}
