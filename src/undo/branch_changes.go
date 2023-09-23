package undo

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/runstate"
	"github.com/git-town/git-town/v9/src/slice"
	"github.com/git-town/git-town/v9/src/steps"
)

type BranchChanges struct {
	LocalAdded    domain.LocalBranchNames
	LocalRemoved  domain.LocalBranchesSHAs
	LocalChanged  domain.LocalBranchChange
	RemoteAdded   domain.RemoteBranchNames
	RemoteRemoved domain.RemoteBranchesSHAs
	RemoteChanged domain.RemoteBranchChange
	OmniRemoved   domain.LocalBranchesSHAs
	// OmniChanges are changes where the local SHA and the remote SHA are identical before the change as well as after the change, and the SHA before and the SHA after are different.
	// Git Town recognizes OmniChanges because only they allow undoing changes made to remote perennial branches.
	// The reason is that perennial branches have protected remote branches, i.e. don't allow force-pushes to their remote branch. One can only do normal pushes.
	// So, to revert a change on a remote perennial branch one needs to perform a revert commit on the local perennial branch,
	// then normal-push (not force-push) that new commit up to the remote branch.
	// This is only possible if the local and remote branches have an identical SHA before as well as after.
	OmniChanged domain.LocalBranchChange // a branch had the same SHA locally and remotely, now it has a new SHA locally and remotely, the local and remote SHA are still equal
	// Inconsistent changes are changes on both local and tracking branch, but where the local and tracking branch don't have the same SHA before or after.
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

func (c BranchChanges) UndoSteps(args StepsArgs) runstate.StepList {
	fmt.Println("111111111111111111111111")
	fmt.Println(c)
	result := runstate.StepList{}
	omniChangedPerennials, omniChangedFeatures := c.OmniChanged.Categorize(args.BranchTypes)

	// revert omni-changed perennial branches
	for branch, change := range omniChangedPerennials {
		if slice.Contains(args.UndoablePerennialCommits, change.After) {
			result.Append(&steps.CheckoutStep{Branch: branch})
			result.Append(&steps.RevertCommitStep{SHA: change.After})
			result.Append(&steps.PushCurrentBranchStep{CurrentBranch: branch, NoPushHook: true})
		}
	}

	// reset omni-changed feature branches
	for branch, change := range omniChangedFeatures {
		result.Append(&steps.CheckoutStep{Branch: branch})
		result.Append(&steps.ResetCurrentBranchToSHAStep{MustHaveSHA: change.After, SetToSHA: change.Before, Hard: true})
		result.Append(&steps.ForcePushBranchStep{Branch: branch, NoPushHook: true})
	}

	// re-create removed omni-branches
	for branch, sha := range c.OmniRemoved {
		result.Append(&steps.CreateBranchStep{Branch: branch, StartingPoint: sha.Location()})
		result.Append(&steps.CreateTrackingBranchStep{Branch: branch, NoPushHook: true})
	}

	// ignore inconsistently changed perennial branches
	// because we can't change the remote and we therefore don't want to reset the local part either
	_, inconsistentChangedFeatures := c.InconsistentlyChanged.Categorize(args.BranchTypes)

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
		result.Append(&steps.DeleteTrackingBranchStep{
			Branch: addedRemoteBranch.LocalBranchName(),
		})
	}

	// re-create remotely removed feature branches
	_, removedFeatureTrackingBranches := c.RemoteRemoved.Categorize(args.BranchTypes)
	for branch, sha := range removedFeatureTrackingBranches {
		result.Append(&steps.CreateRemoteBranchStep{
			Branch:     branch.LocalBranchName(),
			SHA:        sha,
			NoPushHook: true,
		})
	}

	// reset locally changed branches
	for localBranch, change := range c.LocalChanged {
		result.Append(&steps.CheckoutStep{Branch: localBranch})
		result.Append(&steps.ResetCurrentBranchToSHAStep{MustHaveSHA: change.After, SetToSHA: change.Before, Hard: true})
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

	// re-create locally removed branches
	for removedLocalBranch, startingPoint := range c.LocalRemoved {
		result.Append(&steps.CreateBranchStep{
			Branch:        removedLocalBranch,
			StartingPoint: startingPoint.Location(),
		})
	}

	_, remoteFeatureChanges := c.RemoteChanged.Categorize(args.BranchTypes)
	// Ignore remotely changed perennial branches because we can't force-push to them
	// and we would need the local branch to revert commits on them, but we can't change the local branch.

	// reset remotely changed feature branches
	for remoteChangedFeatureBranch, change := range remoteFeatureChanges {
		result.Append(&steps.ResetRemoteBranchToSHAStep{
			Branch:      remoteChangedFeatureBranch,
			MustHaveSHA: change.After,
			SetToSHA:    change.Before,
		})
	}

	result.Append(&steps.CheckoutStep{Branch: args.InitialBranch})
	return result
}

type StepsArgs struct {
	Lineage                  config.Lineage
	BranchTypes              domain.BranchTypes
	InitialBranch            domain.LocalBranchName
	FinalBranch              domain.LocalBranchName
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
