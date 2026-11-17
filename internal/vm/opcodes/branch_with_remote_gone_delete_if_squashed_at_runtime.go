package opcodes

import (
	"github.com/git-town/git-town/v24/internal/config/configdomain"
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/vm/shared"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

// BranchWithRemoteGoneDeleteIfSquashedAtRuntime deletes the given branch if its content is already in its parent,
// for example because the branch was squash-merged at the remote.
// Otherwise it pulls the parent branch into the given branch.
type BranchWithRemoteGoneDeleteIfSquashedAtRuntime struct {
	Branch            gitdomain.LocalBranchName
	CommitsToRemove   Option[gitdomain.SHA]
	InitialParentName Option[gitdomain.LocalBranchName]
	InitialParentSHA  Option[gitdomain.SHA]
	SyncStrategy      configdomain.SyncFeatureStrategy
}

func (self *BranchWithRemoteGoneDeleteIfSquashedAtRuntime) Run(args shared.RunArgs) error {
	// if the content of the shipped branch is in the parent --> delete it
	isContentInParent, err := self.isContentInParent(args)
	if err != nil {
		return err
	}
	if isContentInParent {
		args.PrependOpcodes(
			&CheckoutDescendentOrOther{Branch: self.Branch},
			&BranchLocalDeleteContent{
				BranchToDelete:     self.Branch,
				BranchToRebaseOnto: self.InitialParentName.GetOr(args.Config.Value.ValidatedConfigData.MainBranch),
			},
			&LineageBranchRemove{Branch: self.Branch},
		)
		return nil
	}
	// here the remote ship was incomplete, the local branch has changes that aren't in the parent --> sync the local branch
	switch self.SyncStrategy {
	case configdomain.SyncFeatureStrategyMerge, configdomain.SyncFeatureStrategyCompress:
		args.PrependOpcodes(&SyncFeatureBranchMerge{
			Branch:            self.Branch,
			InitialParentName: self.InitialParentName,
			InitialParentSHA:  self.InitialParentSHA,
			TrackingBranch:    None[gitdomain.RemoteBranchName](),
		})
	case configdomain.SyncFeatureStrategyRebase:
		args.PrependOpcodes(&RebaseAncestorsUntilLocal{
			Branch:          self.Branch,
			CommitsToRemove: self.CommitsToRemove,
		})
	}
	return nil
}

// isContentInParent indicates whether the content of the given branch is already in its parent.
// first checks whether the files changed by the branch have the parent's content,
// then whether the parent contains the squashed changes of the branch as a single commit.
func (self *BranchWithRemoteGoneDeleteIfSquashedAtRuntime) isContentInParent(args shared.RunArgs) (bool, error) {
	parent := self.InitialParentName.GetOr(args.Config.Value.ValidatedConfigData.MainBranch)
	isSameContent, err := args.Git.BranchContentInParent(args.Backend, self.Branch, parent)
	if err != nil || isSameContent {
		return isSameContent, err
	}
	return args.Git.BranchSquashAppliedInParent(args.Backend, self.Branch, parent)
}
