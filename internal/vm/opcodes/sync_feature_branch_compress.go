package opcodes

import (
	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/vm/shared"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// SyncFeatureBranchCompress expands to all opcodes needed to sync a feature branch using the "compress" sync strategy.
type SyncFeatureBranchCompress struct {
	CommitMessage           Option[gitdomain.CommitMessage]
	CurrentBranch           gitdomain.LocalBranchName
	InitialParentName       Option[gitdomain.LocalBranchName]
	InitialParentSHA        Option[gitdomain.SHA]
	Offline                 configdomain.Offline
	TrackingBranch          Option[gitdomain.RemoteBranchName]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *SyncFeatureBranchCompress) Run(args shared.RunArgs) error {
	opcodes := []shared.Opcode{}
	commitsInBranch := gitdomain.Commits{}
	if parentLocalName, hasParent := args.Config.Value.NormalConfig.Lineage.Parent(self.CurrentBranch).Get(); hasParent {
		parentName := determineParentBranchName(parentLocalName, args.BranchInfos, args.Config.Value.NormalConfig.DevRemote)
		inSyncWithParent, err := args.Git.BranchInSyncWithParent(args.Backend, self.CurrentBranch, parentName)
		if err != nil {
			return err
		}
		if !inSyncWithParent {
			opcodes = append(opcodes, &SyncFeatureBranchMerge{
				Branch:            self.CurrentBranch,
				InitialParentName: self.InitialParentName,
				InitialParentSHA:  self.InitialParentSHA,
			})
		}
		commitsInBranch, err = args.Git.CommitsInFeatureBranch(args.Backend, self.CurrentBranch, parentName)
		if err != nil {
			return err
		}
	}
	if trackingBranch, hasTrackingBranch := self.TrackingBranch.Get(); hasTrackingBranch {
		inSyncWithTracking, err := args.Git.BranchInSyncWithTracking(args.Backend, self.CurrentBranch, args.Config.Value.NormalConfig.DevRemote)
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
		if self.Offline.IsFalse() && self.TrackingBranch.IsSome() {
			opcodes = append(opcodes, &PushCurrentBranchForceIfNeeded{CurrentBranch: self.CurrentBranch, ForceIfIncludes: false})
		}
	}
	args.PrependOpcodes(opcodes...)
	return nil
}

func determineParentBranchName(parentLocalName gitdomain.LocalBranchName, branchInfosOpt Option[gitdomain.BranchInfos], devRemote gitdomain.Remote) gitdomain.BranchName {
	if branchInfos, hasBranchInfos := branchInfosOpt.Get(); hasBranchInfos {
		if parentInfo, hasParentInfo := branchInfos.FindByLocalName(parentLocalName).Get(); hasParentInfo {
			return parentInfo.GetLocalOrRemoteName()
		}
		parentRemoteName := parentLocalName.AtRemote(devRemote)
		if _, hasParentInfo := branchInfos.FindByRemoteName(parentRemoteName).Get(); hasParentInfo {
			return parentRemoteName.BranchName()
		}
	}
	return parentLocalName.BranchName()
}
