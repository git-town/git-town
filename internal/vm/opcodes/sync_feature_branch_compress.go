package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	"github.com/git-town/git-town/v22/internal/vm/shared"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// SyncFeatureBranchCompress expands to all opcodes needed to sync a feature branch using the "compress" sync strategy.
type SyncFeatureBranchCompress struct {
	CommitMessage     Option[gitdomain.CommitMessage]
	CurrentBranch     gitdomain.LocalBranchName
	InitialParentName Option[gitdomain.LocalBranchName]
	InitialParentSHA  Option[gitdomain.SHA]
	Offline           configdomain.Offline
	PushBranches      configdomain.PushBranches
	TrackingBranch    Option[gitdomain.RemoteBranchName]
}

func (self *SyncFeatureBranchCompress) Run(args shared.RunArgs) error {
	opcodes := []shared.Opcode{}
	commitsInBranch := gitdomain.Commits{}
	if parentLocalName, hasParent := args.Config.Value.NormalConfig.Lineage.Parent(self.CurrentBranch).Get(); hasParent {
		parentName, hasParent := determineParentBranchName(parentLocalName, args.BranchInfos).Get()
		if !hasParent {
			return fmt.Errorf(messages.BranchInfoNotFound, parentLocalName)
		}
		inSyncWithParent, err := args.Git.BranchInSyncWithParent(args.Backend, self.CurrentBranch, parentName)
		if err != nil {
			return err
		}
		parentIsPerennial := args.Config.Value.IsMainOrPerennialBranch(parentLocalName)
		skipParent := args.Config.Value.NormalConfig.Detached.ShouldWorkDetached() && parentIsPerennial
		if !inSyncWithParent && !skipParent {
			opcodes = append(opcodes, &SyncFeatureBranchMerge{
				Branch:            self.CurrentBranch,
				InitialParentName: self.InitialParentName,
				InitialParentSHA:  self.InitialParentSHA,
				// We must sync with the tracking branch separately below,
				// because this only runs if we aren't in sync with the parent.
				TrackingBranch: None[gitdomain.RemoteBranchName](),
			})
		}
		commitsInBranch, err = args.Git.CommitsInFeatureBranch(args.Backend, self.CurrentBranch, parentName)
		if err != nil {
			return err
		}
	}
	if trackingBranch, hasTrackingBranch := self.TrackingBranch.Get(); hasTrackingBranch {
		inSyncWithTracking, err := args.Git.BranchInSyncWithTracking(args.Backend, self.CurrentBranch, trackingBranch)
		if err != nil {
			return err
		}
		if !inSyncWithTracking {
			opcodes = append(opcodes,
				&MergeIntoCurrentBranch{BranchToMerge: trackingBranch.BranchName()},
			)
		}
	}
	commitMessage, hasCommitMessage := self.CommitMessage.Get()
	if hasCommitMessage && (len(opcodes) > 0 || len(commitsInBranch) > 1) {
		opcodes = append(opcodes,
			&BranchCurrentResetToParent{
				CurrentBranch: self.CurrentBranch,
			},
			&CommitWithMessage{
				AuthorOverride: None[gitdomain.Author](),
				CommitHook:     configdomain.CommitHookEnabled,
				Message:        commitMessage,
			},
		)
		trackingBranch, hasTrackingBranch := self.TrackingBranch.Get()
		if self.Offline.IsOnline() && hasTrackingBranch && self.PushBranches.ShouldPush() {
			opcodes = append(opcodes, &PushCurrentBranchForceIfNeeded{
				CurrentBranch:   self.CurrentBranch,
				ForceIfIncludes: false,
				TrackingBranch:  trackingBranch,
			})
		}
	}
	args.PrependOpcodes(opcodes...)
	return nil
}

func determineParentBranchName(parentLocalName gitdomain.LocalBranchName, branchInfos gitdomain.BranchInfos) Option[gitdomain.BranchName] {
	if localInfo, parentIsLocal := branchInfos.FindByLocalName(parentLocalName).Get(); parentIsLocal {
		return Some(localInfo.GetLocalOrRemoteName())
	}
	if trackingInfo, parentIsRemote := branchInfos.FindRemoteNameMatchingLocal(parentLocalName).Get(); parentIsRemote {
		return Some(trackingInfo.GetLocalOrRemoteName())
	}
	return None[gitdomain.BranchName]()
}
