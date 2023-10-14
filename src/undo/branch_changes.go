package undo

import (
	"strings"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/gohacks/slice"
	"github.com/git-town/git-town/v9/src/step"
	"github.com/git-town/git-town/v9/src/vm/program"
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
		LocalRemoved:          domain.LocalBranchesSHAs{},
		LocalChanged:          domain.LocalBranchChange{},
		RemoteAdded:           domain.RemoteBranchNames{},
		RemoteRemoved:         map[domain.RemoteBranchName]domain.SHA{},
		RemoteChanged:         map[domain.RemoteBranchName]domain.Change[domain.SHA]{},
		OmniRemoved:           domain.LocalBranchesSHAs{},
		OmniChanged:           domain.LocalBranchChange{},
		InconsistentlyChanged: domain.InconsistentChanges{},
	}
}

// UndoSteps provides the steps to undo the changes described by this BranchChanges instance.
func (bcs BranchChanges) UndoSteps(args StepsArgs) program.List {
	result := program.List{}
	omniChangedPerennials, omniChangedFeatures := bcs.OmniChanged.Categorize(args.BranchTypes)

	// revert omni-changed perennial branches
	for _, branch := range omniChangedPerennials.BranchNames() {
		change := omniChangedPerennials[branch]
		if slice.Contains(args.UndoablePerennialCommits, change.After) {
			result.Add(&step.Checkout{Branch: branch})
			result.Add(&step.RevertCommit{SHA: change.After})
			result.Add(&step.PushCurrentBranch{CurrentBranch: branch, NoPushHook: args.NoPushHook})
		}
	}

	// reset omni-changed feature branches
	for _, branch := range omniChangedFeatures.BranchNames() {
		change := omniChangedFeatures[branch]
		result.Add(&step.Checkout{Branch: branch})
		result.Add(&step.ResetCurrentBranchToSHA{MustHaveSHA: change.After, SetToSHA: change.Before, Hard: true})
		result.Add(&step.ForcePushCurrentBranch{NoPushHook: args.NoPushHook})
	}

	// re-create removed omni-branches
	for _, branch := range bcs.OmniRemoved.BranchNames() {
		sha := bcs.OmniRemoved[branch]
		result.Add(&step.CreateBranch{Branch: branch, StartingPoint: sha.Location()})
		result.Add(&step.CreateTrackingBranch{Branch: branch, NoPushHook: args.NoPushHook})
	}

	inconsistentlyChangedPerennials, inconsistentChangedFeatures := bcs.InconsistentlyChanged.Categorize(args.BranchTypes)

	// reset inconsintently changed perennial branches
	for _, inconsistentlyChangedPerennial := range inconsistentlyChangedPerennials {
		if inconsistentlyChangedPerennial.After.LocalSHA == inconsistentlyChangedPerennial.After.RemoteSHA {
			if slice.Contains(args.UndoablePerennialCommits, inconsistentlyChangedPerennial.After.LocalSHA) {
				result.Add(&step.Checkout{Branch: inconsistentlyChangedPerennial.Before.LocalName})
				result.Add(&step.RevertCommit{SHA: inconsistentlyChangedPerennial.After.LocalSHA})
				result.Add(&step.PushCurrentBranch{CurrentBranch: inconsistentlyChangedPerennial.After.LocalName, NoPushHook: args.NoPushHook})
			}
		}
	}

	// reset inconsintently changed feature branches
	for _, inconsistentChange := range inconsistentChangedFeatures {
		result.Add(&step.Checkout{Branch: inconsistentChange.Before.LocalName})
		result.Add(&step.ResetCurrentBranchToSHA{
			MustHaveSHA: inconsistentChange.After.LocalSHA,
			SetToSHA:    inconsistentChange.Before.LocalSHA,
			Hard:        true,
		})
		result.Add(&step.ResetRemoteBranchToSHA{
			Branch:      inconsistentChange.Before.RemoteName,
			MustHaveSHA: inconsistentChange.After.RemoteSHA,
			SetToSHA:    inconsistentChange.Before.RemoteSHA,
		})
	}

	// remove remotely added branches
	for _, addedRemoteBranch := range bcs.RemoteAdded {
		if addedRemoteBranch.Remote() != domain.UpstreamRemote {
			result.Add(&step.DeleteTrackingBranch{
				Branch: addedRemoteBranch,
			})
		}
	}

	// re-create remotely removed feature branches
	_, removedFeatureTrackingBranches := bcs.RemoteRemoved.Categorize(args.BranchTypes)
	for _, branch := range removedFeatureTrackingBranches.BranchNames() {
		sha := removedFeatureTrackingBranches[branch]
		result.Add(&step.CreateRemoteBranch{
			Branch:     branch.LocalBranchName(),
			SHA:        sha,
			NoPushHook: args.NoPushHook,
		})
	}

	// reset locally changed branches
	for _, localBranch := range bcs.LocalChanged.BranchNames() {
		change := bcs.LocalChanged[localBranch]
		result.Add(&step.Checkout{Branch: localBranch})
		result.Add(&step.ResetCurrentBranchToSHA{MustHaveSHA: change.After, SetToSHA: change.Before, Hard: true})
	}

	// re-create locally removed branches
	for _, removedLocalBranch := range bcs.LocalRemoved.BranchNames() {
		startingPoint := bcs.LocalRemoved[removedLocalBranch]
		result.Add(&step.CreateBranch{
			Branch:        removedLocalBranch,
			StartingPoint: startingPoint.Location(),
		})
	}

	// remove locally added branches
	for _, addedLocalBranch := range bcs.LocalAdded {
		if args.FinalBranch == addedLocalBranch {
			result.Add(&step.Checkout{Branch: args.InitialBranch})
		}
		result.Add(&step.DeleteLocalBranch{
			Branch: addedLocalBranch,
			Force:  true,
		})
	}

	// Ignore remotely changed perennial branches because we can't force-push to them
	// and we would need the local branch to revert commits on them, but we can't change the local branch.

	// reset remotely changed feature branches
	_, remoteFeatureChanges := bcs.RemoteChanged.Categorize(args.BranchTypes)
	for _, remoteChangedFeatureBranch := range remoteFeatureChanges.BranchNames() {
		change := remoteFeatureChanges[remoteChangedFeatureBranch]
		result.Add(&step.ResetRemoteBranchToSHA{
			Branch:      remoteChangedFeatureBranch,
			MustHaveSHA: change.After,
			SetToSHA:    change.Before,
		})
	}

	// This must be a CheckoutIfExistsStep because this branch might not exist
	// when a Git Town command fails, stores this undo step, then gets continued and deletes this branch.
	result.Add(&step.CheckoutIfExists{Branch: args.InitialBranch})
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

func (bcs BranchChanges) String() string {
	s := strings.Builder{}
	s.WriteString("BranchChanges {")
	s.WriteString("\n  LocalAdded: ")
	s.WriteString(strings.Join(bcs.LocalAdded.Strings(), ", "))
	s.WriteString("\n  LocalRemoved: ")
	s.WriteString(strings.Join(bcs.LocalRemoved.BranchNames().Strings(), ", "))
	s.WriteString("\n  LocalChanged: ")
	s.WriteString(strings.Join(bcs.LocalChanged.BranchNames().Strings(), ", "))
	s.WriteString("\n  RemoteAdded: ")
	s.WriteString(strings.Join(bcs.RemoteAdded.Strings(), ", "))
	s.WriteString("\n  RemoteRemoved: ")
	s.WriteString(strings.Join(bcs.RemoteRemoved.BranchNames().Strings(), ", "))
	s.WriteString("\n  RemoteChanged: ")
	s.WriteString(strings.Join(bcs.RemoteChanged.BranchNames().Strings(), ", "))
	s.WriteString("\n  OmniRemoved: ")
	s.WriteString(strings.Join(bcs.OmniRemoved.BranchNames().Strings(), ", "))
	s.WriteString("\n  OmniChanged: ")
	s.WriteString(strings.Join(bcs.OmniChanged.BranchNames().Strings(), ", "))
	s.WriteString("\n  InconsistentlyChanged: ")
	s.WriteString(strings.Join(bcs.InconsistentlyChanged.BranchNames().Strings(), ", "))
	s.WriteString("\n")
	return s.String()
}
