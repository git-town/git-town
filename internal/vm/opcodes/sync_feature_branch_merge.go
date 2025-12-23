package opcodes

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/vm/shared"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// SyncFeatureBranchMerge merges the parent branches of the given branch until a local parent is found.
type SyncFeatureBranchMerge struct {
	Branch            gitdomain.LocalBranchName
	InitialParentName Option[gitdomain.LocalBranchName]
	InitialParentSHA  Option[gitdomain.SHA]
	TrackingBranch    Option[gitdomain.RemoteBranchName]
}

func (self *SyncFeatureBranchMerge) Run(args shared.RunArgs) error {
	program := []shared.Opcode{}
	branch := self.Branch
	for {
		parent, hasParent := args.Config.Value.NormalConfig.Lineage.Parent(branch).Get()
		if !hasParent {
			break
		}
		parentIsPerennial := args.Config.Value.IsMainOrPerennialBranch(parent)
		if args.Config.Value.NormalConfig.Detached.ShouldWorkDetached() && parentIsPerennial {
			break
		}
		if parentBranchInfo, hasParentInfo := args.BranchInfos.FindLocalOrRemote(parent, args.Config.Value.NormalConfig.DevRemote).Get(); hasParentInfo {
			parentIsLocal := parentBranchInfo.LocalName.IsSome()
			if parentIsLocal {
				isInSync, err := args.Git.BranchInSyncWithParent(args.Backend, self.Branch, parent.BranchName())
				if err != nil {
					return err
				}
				if !isInSync {
					program = append(program, &MergeParentResolvePhantomConflicts{
						CurrentBranch:     self.Branch,
						CurrentParent:     parent.BranchName(),
						InitialParentName: self.InitialParentName,
						InitialParentSHA:  self.InitialParentSHA,
					})
				}
				break
			}
			// here the parent isn't local --> sync with its tracking branch if it exists, then try again with the grandparent until we find a local ancestor
			if parentTrackingBranch, parentHasTrackingBranch := parentBranchInfo.RemoteName.Get(); parentHasTrackingBranch {
				isInSync, err := args.Git.BranchInSyncWithParent(args.Backend, self.Branch, parentTrackingBranch.BranchName())
				if err != nil {
					return err
				}
				if !isInSync {
					program = append(program, &MergeParentResolvePhantomConflicts{
						CurrentBranch:     self.Branch,
						CurrentParent:     parentTrackingBranch.BranchName(),
						InitialParentName: self.InitialParentName,
						InitialParentSHA:  self.InitialParentSHA,
					})
				}
			}
		}
		branch = parent
	}
	if trackingBranch, hasTrackingBranch := self.TrackingBranch.Get(); hasTrackingBranch {
		isInSync, err := args.Git.BranchInSyncWithTracking(args.Backend, self.Branch, trackingBranch)
		if err != nil {
			return err
		}
		if !isInSync {
			program = append(program, &MergeIntoCurrentBranch{BranchToMerge: trackingBranch.BranchName()})
		}
	}
	args.PrependOpcodes(program...)
	return nil
}
