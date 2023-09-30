package undo

import (
	"strings"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/slice"
	"github.com/git-town/git-town/v9/src/steps"
)

// BranchChanges describes the changes made to the branches in a Git repo.
// Various types of changes are distinguished.
type BranchChanges struct {
	LocalAdded    domain.LocalBranchNames
	LocalRemoved  domain.LocalBranchesSHAs
	LocalChanged  domain.LocalBranchChange
	RemoteAdded   domain.RemoteBranchNames
	RemoteRemoved domain.RemoteBranchesSHAs
	RemoteChanged domain.RemoteBranchChange
	// OmniRemoved is when a branch that has the same SHA on its local and tracking branch gets removed.
	OmniRemoved domain.LocalBranchesSHAs
	// OmniChanges are changes where the local SHA and the remote SHA are identical before the change as well as after the change,
	OmniChanged domain.LocalBranchChange // a branch had the same SHA locally and remotely, now it has a new SHA locally and remotely, the local and remote SHA are still equal
	// Inconsistent changes are changes on both local and tracking branch, but where the local and tracking branch
	// don't have the same SHA before or after.
	// These changes cannot be undone for perennial branches because there is no way to reset the remote branch to the SHA it had before.
	InconsistentlyChanged domain.InconsistentChanges
}

// EmptyBranchChanges provides a properly initialized empty Changes instance.
func EmptyBranchChanges() BranchChanges {
	return BranchChanges{
		LocalAdded:            domain.LocalBranchNames{},
		LocalRemoved:          map[domain.LocalBranchName]domain.SHA{},
		LocalChanged:          domain.LocalBranchChange{},
		RemoteAdded:           []domain.RemoteBranchName{},
		RemoteRemoved:         map[domain.RemoteBranchName]domain.SHA{},
		RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
		OmniRemoved:           map[domain.LocalBranchName]domain.SHA{},
		OmniChanged:           domain.LocalBranchChange{},
		InconsistentlyChanged: domain.InconsistentChanges{},
	}
}

// UndoSteps provides the steps to undo the changes described by this BranchChanges instance.
func (c BranchChanges) UndoSteps(args StepsArgs) runstate.StepList {
	result := runstate.StepList{}
	omniChangedPerennials, omniChangedFeatures := c.OmniChanged.Categorize(args.BranchTypes)

	// revert omni-changed perennial branches
	for _, branch := range omniChangedPerennials.BranchNames() {
		change := omniChangedPerennials[branch]
		if slice.Contains(args.UndoablePerennialCommits, change.After) {
			result.Append(&steps.CheckoutStep{Branch: branch})
			result.Append(&steps.RevertCommitStep{SHA: change.After})
			result.Append(&steps.PushCurrentBranchStep{CurrentBranch: branch, NoPushHook: args.NoPushHook})
		}
	}

	// reset omni-changed feature branches
	for _, branch := range omniChangedFeatures.BranchNames() {
		change := omniChangedFeatures[branch]
		result.Append(&steps.CheckoutStep{Branch: branch})
		result.Append(&steps.ResetCurrentBranchToSHAStep{MustHaveSHA: change.After, SetToSHA: change.Before, Hard: true})
		result.Append(&steps.ForcePushCurrentBranchStep{NoPushHook: args.NoPushHook})
	}

	// re-create removed omni-branches
	for _, branch := range c.OmniRemoved.BranchNames() {
		sha := c.OmniRemoved[branch]
		result.Append(&steps.CreateBranchStep{Branch: branch, StartingPoint: sha.Location()})
		result.Append(&steps.CreateTrackingBranchStep{Branch: branch, NoPushHook: args.NoPushHook})
	}

	inconsistentlyChangedPerennials, inconsistentChangedFeatures := c.InconsistentlyChanged.Categorize(args.BranchTypes)

	// reset inconsintently changed perennial branches
	for _, inconsistentlyChangedPerennial := range inconsistentlyChangedPerennials {
		if inconsistentlyChangedPerennial.After.LocalSHA == inconsistentlyChangedPerennial.After.RemoteSHA {
			if slice.Contains(args.UndoablePerennialCommits, inconsistentlyChangedPerennial.After.LocalSHA) {
				result.Append(&steps.CheckoutStep{Branch: inconsistentlyChangedPerennial.Before.LocalName})
				result.Append(&steps.RevertCommitStep{SHA: inconsistentlyChangedPerennial.After.LocalSHA})
				result.Append(&steps.PushCurrentBranchStep{CurrentBranch: inconsistentlyChangedPerennial.After.LocalName, NoPushHook: args.NoPushHook})
			}
		}
	}

	// reset inconsintently changed feature branches
	for _, inconsistentChange := range inconsistentChangedFeatures {
		result.Append(&steps.CheckoutStep{Branch: inconsistentChange.Before.LocalName})
		result.Append(&steps.ResetCurrentBranchToSHAStep{
			MustHaveSHA: inconsistentChange.After.LocalSHA,
			SetToSHA:    inconsistentChange.Before.LocalSHA,
			Hard:        true,
		})
		result.Append(&steps.ResetRemoteBranchToSHAStep{
			Branch:      inconsistentChange.Before.RemoteName,
			MustHaveSHA: inconsistentChange.After.RemoteSHA,
			SetToSHA:    inconsistentChange.Before.RemoteSHA,
		})
	}

	// remove remotely added branches
	for _, addedRemoteBranch := range c.RemoteAdded {
		if addedRemoteBranch.Remote() != domain.UpstreamRemote {
			result.Append(&steps.DeleteTrackingBranchStep{
				Branch: addedRemoteBranch,
			})
		}
	}

	// re-create remotely removed feature branches
	_, removedFeatureTrackingBranches := c.RemoteRemoved.Categorize(args.BranchTypes)
	for _, branch := range removedFeatureTrackingBranches.BranchNames() {
		sha := removedFeatureTrackingBranches[branch]
		result.Append(&steps.CreateRemoteBranchStep{
			Branch:     branch.LocalBranchName(),
			SHA:        sha,
			NoPushHook: args.NoPushHook,
		})
	}

	// reset locally changed branches
	for _, localBranch := range c.LocalChanged.BranchNames() {
		change := c.LocalChanged[localBranch]
		result.Append(&steps.CheckoutStep{Branch: localBranch})
		result.Append(&steps.ResetCurrentBranchToSHAStep{MustHaveSHA: change.After, SetToSHA: change.Before, Hard: true})
	}

	// re-create locally removed branches
	for _, removedLocalBranch := range c.LocalRemoved.BranchNames() {
		startingPoint := c.LocalRemoved[removedLocalBranch]
		result.Append(&steps.CreateBranchStep{
			Branch:        removedLocalBranch,
			StartingPoint: startingPoint.Location(),
		})
	}

	// remove locally added branches
	for _, addedLocalBranch := range c.LocalAdded {
		if args.FinalBranch == addedLocalBranch {
			result.Append(&steps.CheckoutStep{Branch: args.InitialBranch})
		}
		result.Append(&steps.DeleteLocalBranchStep{
			Branch: addedLocalBranch,
			Parent: args.Lineage.Parent(addedLocalBranch).Location(),
			Force:  true,
		})
	}

	// Ignore remotely changed perennial branches because we can't force-push to them
	// and we would need the local branch to revert commits on them, but we can't change the local branch.

	// reset remotely changed feature branches
	_, remoteFeatureChanges := c.RemoteChanged.Categorize(args.BranchTypes)
	for _, remoteChangedFeatureBranch := range remoteFeatureChanges.BranchNames() {
		change := remoteFeatureChanges[remoteChangedFeatureBranch]
		result.Append(&steps.ResetRemoteBranchToSHAStep{
			Branch:      remoteChangedFeatureBranch,
			MustHaveSHA: change.After,
			SetToSHA:    change.Before,
		})
	}

	// This must be a CheckoutIfExistsStep because this branch might not exist
	// when a Git Town command fails, stores this undo step, then gets continued and deletes this branch.
	result.Append(&steps.CheckoutIfExistsStep{Branch: args.InitialBranch})
	return result
}

type StepsArgs struct {
	Lineage                  config.Lineage
	BranchTypes              domain.BranchTypes
	InitialBranch            domain.LocalBranchName
	FinalBranch              domain.LocalBranchName
	NoPushHook               bool
	UndoablePerennialCommits []domain.SHA
}

func (c BranchChanges) String() string {
	s := strings.Builder{}
	s.WriteString("BranchChanges {")
	s.WriteString("\n  LocalAdded: ")
	s.WriteString(strings.Join(c.LocalAdded.Strings(), ", "))
	s.WriteString("\n  LocalRemoved: ")
	s.WriteString(strings.Join(c.LocalRemoved.BranchNames().Strings(), ", "))
	s.WriteString("\n  LocalChanged: ")
	s.WriteString(strings.Join(c.LocalChanged.BranchNames().Strings(), ", "))
	s.WriteString("\n  RemoteAdded: ")
	s.WriteString(strings.Join(c.RemoteAdded.Strings(), ", "))
	s.WriteString("\n  RemoteRemoved: ")
	s.WriteString(strings.Join(c.RemoteRemoved.BranchNames().Strings(), ", "))
	s.WriteString("\n  RemoteChanged: ")
	s.WriteString(strings.Join(c.RemoteChanged.BranchNames().Strings(), ", "))
	s.WriteString("\n  OmniRemoved: ")
	s.WriteString(strings.Join(c.OmniRemoved.BranchNames().Strings(), ", "))
	s.WriteString("\n  OmniChanged: ")
	s.WriteString(strings.Join(c.OmniChanged.BranchNames().Strings(), ", "))
	s.WriteString("\n  InconsistentlyChanged: ")
	s.WriteString(strings.Join(c.InconsistentlyChanged.BranchNames().Strings(), ", "))
	s.WriteString("\n")
	return s.String()
}
