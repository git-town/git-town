package opcodes

import (
	"fmt"

	"github.com/git-town/git-town/v20/internal/config/configdomain"
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/vm/shared"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// SyncFeatureBranchCompress expands to all opcodes needed to sync a feature branch using the "compress" sync strategy.
type SyncFeatureBranchCompress struct {
	CommitMessage           Option[gitdomain.CommitMessage]
	CurrentBranch           gitdomain.LocalBranchName
	DevRemote               gitdomain.Remote
	Offline                 configdomain.Offline
	ParentName              Option[gitdomain.LocalBranchName]
	ParentSHA               Option[gitdomain.SHA]
	TrackingBranch          Option[gitdomain.RemoteBranchName]
	undeclaredOpcodeMethods `exhaustruct:"optional"`
}

func (self *SyncFeatureBranchCompress) Run(args shared.RunArgs) error {
	opcodes := []shared.Opcode{}
	commitsInBranch := gitdomain.Commits{}
	if parentName, hasParent := self.ParentName.Get(); hasParent {
		inSyncWithParent, err := args.Git.BranchInSyncWithParent(args.Backend, self.CurrentBranch, parentName)
		if err != nil {
			return err
		}
		if !inSyncWithParent {
			opcodes = append(opcodes, &MergeParentsUntilLocal{
				Branch:             self.CurrentBranch,
				OriginalParentName: self.ParentName,
				OriginalParentSHA:  self.ParentSHA,
			})
		}
		commitsInBranch, err = args.Git.CommitsInFeatureBranch(args.Backend, self.CurrentBranch, parentName)
	}
	if trackingBranch, hasTrackingBranch := self.TrackingBranch.Get(); hasTrackingBranch {
		inSyncWithTracking, err := args.Git.BranchInSyncWithTracking(args.Backend, self.CurrentBranch, self.DevRemote)
		if err != nil {
			return err
		}
		if !inSyncWithTracking {
			opcodes = append(opcodes,
				&MergeIntoCurrentBranch{BranchToMerge: trackingBranch.BranchName()},
			)
		}
	}
	fmt.Println("1111111111111111111111111111111111")
	fmt.Println(opcodes)
	fmt.Println(commitsInBranch)
	if len(opcodes) > 0 || len(commitsInBranch) > 1 {
		firstCommitMessage := self.CommitMessage.GetOrElse("compressed commit")
		opcodes = append(opcodes,
			&BranchCurrentResetToParent{
				CurrentBranch: self.CurrentBranch,
			},
			&CommitWithMessage{
				AuthorOverride: None[gitdomain.Author](),
				CommitHook:     configdomain.CommitHookEnabled,
				Message:        firstCommitMessage,
			},
		)
		if self.Offline.IsFalse() {
			opcodes = append(opcodes, &PushCurrentBranchForceIfNeeded{CurrentBranch: self.CurrentBranch, ForceIfIncludes: false})
		}
	}
	args.PrependOpcodes(opcodes...)
	return nil
}
