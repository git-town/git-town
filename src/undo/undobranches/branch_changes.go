package undobranches

import (
	"strings"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/gohacks/slice"
	"github.com/git-town/git-town/v11/src/undo/undodomain"
	"github.com/git-town/git-town/v11/src/vm/opcode"
	"github.com/git-town/git-town/v11/src/vm/program"
)

// BranchChanges describes the changes made to the branches in a Git repo.
// Various types of changes are distinguished.
type BranchChanges struct {
	LocalAdded    gitdomain.LocalBranchNames
	LocalRemoved  LocalBranchesSHAs
	LocalChanged  LocalBranchChange
	RemoteAdded   gitdomain.RemoteBranchNames
	RemoteRemoved RemoteBranchesSHAs
	RemoteChanged RemoteBranchChange
	// OmniRemoved is when a branch that has the same SHA on its local and tracking branch gets removed.
	OmniRemoved LocalBranchesSHAs
	// OmniChanges are changes where the local SHA and the remote SHA are identical before the change as well as after the change,
	OmniChanged LocalBranchChange // a branch had the same SHA locally and remotely, now it has a new SHA locally and remotely, the local and remote SHA are still equal
	// Inconsistent changes are changes on both local and tracking branch, but where the local and tracking branch
	// don't have the same SHA before or after.
	// These changes cannot be undone for perennial branches because there is no way to reset the remote branch to the SHA it had before.
	InconsistentlyChanged undodomain.InconsistentChanges
}

// EmptyBranchChanges provides a properly initialized empty Changes instance.
func EmptyBranchChanges() BranchChanges {
	return BranchChanges{
		LocalAdded:            gitdomain.LocalBranchNames{},
		LocalRemoved:          LocalBranchesSHAs{},
		LocalChanged:          LocalBranchChange{},
		RemoteAdded:           gitdomain.RemoteBranchNames{},
		RemoteRemoved:         map[gitdomain.RemoteBranchName]gitdomain.SHA{},
		RemoteChanged:         map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
		OmniRemoved:           LocalBranchesSHAs{},
		OmniChanged:           LocalBranchChange{},
		InconsistentlyChanged: undodomain.InconsistentChanges{},
	}
}

func (self BranchChanges) String() string {
	s := strings.Builder{}
	s.WriteString("BranchChanges {")
	s.WriteString("\n  LocalAdded: ")
	s.WriteString(strings.Join(self.LocalAdded.Strings(), ", "))
	s.WriteString("\n  LocalRemoved: ")
	s.WriteString(strings.Join(self.LocalRemoved.BranchNames().Strings(), ", "))
	s.WriteString("\n  LocalChanged: ")
	s.WriteString(strings.Join(self.LocalChanged.BranchNames().Strings(), ", "))
	s.WriteString("\n  RemoteAdded: ")
	s.WriteString(strings.Join(self.RemoteAdded.Strings(), ", "))
	s.WriteString("\n  RemoteRemoved: ")
	s.WriteString(strings.Join(self.RemoteRemoved.BranchNames().Strings(), ", "))
	s.WriteString("\n  RemoteChanged: ")
	s.WriteString(strings.Join(self.RemoteChanged.BranchNames().Strings(), ", "))
	s.WriteString("\n  OmniRemoved: ")
	s.WriteString(strings.Join(self.OmniRemoved.BranchNames().Strings(), ", "))
	s.WriteString("\n  OmniChanged: ")
	s.WriteString(strings.Join(self.OmniChanged.BranchNames().Strings(), ", "))
	s.WriteString("\n  InconsistentlyChanged: ")
	s.WriteString(strings.Join(self.InconsistentlyChanged.BranchNames().Strings(), ", "))
	s.WriteString("\n")
	return s.String()
}

// UndoProgram provides the steps to undo the changes described by this BranchChanges instance.
func (self BranchChanges) UndoProgram(args BranchChangesUndoProgramArgs) program.Program {
	result := program.Program{}
	omniChangedPerennials, omniChangedFeatures := CategorizeLocalBranchChange(self.OmniChanged, args.BranchTypes)

	// revert omni-changed perennial branches
	for _, branch := range omniChangedPerennials.BranchNames() {
		change := omniChangedPerennials[branch]
		if slice.Contains(args.UndoablePerennialCommits, change.After) {
			result.Add(&opcode.Checkout{Branch: branch})
			result.Add(&opcode.RevertCommit{SHA: change.After})
			result.Add(&opcode.PushCurrentBranch{CurrentBranch: branch})
		}
	}

	// reset omni-changed feature branches
	for _, branch := range omniChangedFeatures.BranchNames() {
		change := omniChangedFeatures[branch]
		result.Add(&opcode.Checkout{Branch: branch})
		result.Add(&opcode.ResetCurrentBranchToSHA{MustHaveSHA: change.After, SetToSHA: change.Before, Hard: true})
		result.Add(&opcode.ForcePushCurrentBranch{})
	}

	// re-create removed omni-branches
	for _, branch := range self.OmniRemoved.BranchNames() {
		sha := self.OmniRemoved[branch]
		result.Add(&opcode.CreateBranch{Branch: branch, StartingPoint: sha.Location()})
		result.Add(&opcode.CreateTrackingBranch{Branch: branch})
	}

	inconsistentlyChangedPerennials, inconsistentChangedFeatures := CategorizeInconsistentChanges(self.InconsistentlyChanged, args.BranchTypes)

	// reset inconsintently changed perennial branches
	for _, inconsistentlyChangedPerennial := range inconsistentlyChangedPerennials {
		if inconsistentlyChangedPerennial.After.LocalSHA == inconsistentlyChangedPerennial.After.RemoteSHA {
			if slice.Contains(args.UndoablePerennialCommits, inconsistentlyChangedPerennial.After.LocalSHA) {
				result.Add(&opcode.Checkout{Branch: inconsistentlyChangedPerennial.Before.LocalName})
				result.Add(&opcode.RevertCommit{SHA: inconsistentlyChangedPerennial.After.LocalSHA})
				result.Add(&opcode.PushCurrentBranch{CurrentBranch: inconsistentlyChangedPerennial.After.LocalName})
			}
		}
	}

	// reset inconsintently changed feature branches
	for _, inconsistentChange := range inconsistentChangedFeatures {
		result.Add(&opcode.Checkout{Branch: inconsistentChange.Before.LocalName})
		result.Add(&opcode.ResetCurrentBranchToSHA{
			MustHaveSHA: inconsistentChange.After.LocalSHA,
			SetToSHA:    inconsistentChange.Before.LocalSHA,
			Hard:        true,
		})
		result.Add(&opcode.ResetRemoteBranchToSHA{
			Branch:      inconsistentChange.Before.RemoteName,
			MustHaveSHA: inconsistentChange.After.RemoteSHA,
			SetToSHA:    inconsistentChange.Before.RemoteSHA,
		})
	}

	// remove remotely added branches
	for _, addedRemoteBranch := range self.RemoteAdded {
		if addedRemoteBranch.Remote() != gitdomain.UpstreamRemote {
			result.Add(&opcode.DeleteTrackingBranch{
				Branch: addedRemoteBranch,
			})
		}
	}

	// re-create remotely removed feature branches
	_, removedFeatureTrackingBranches := CategorizeRemoteBranchesSHAs(self.RemoteRemoved, args.BranchTypes)
	for _, branch := range removedFeatureTrackingBranches.BranchNames() {
		sha := removedFeatureTrackingBranches[branch]
		result.Add(&opcode.CreateRemoteBranch{
			Branch: branch.LocalBranchName(),
			SHA:    sha,
		})
	}

	// reset locally changed branches
	for _, localBranch := range self.LocalChanged.BranchNames() {
		change := self.LocalChanged[localBranch]
		result.Add(&opcode.Checkout{Branch: localBranch})
		result.Add(&opcode.ResetCurrentBranchToSHA{MustHaveSHA: change.After, SetToSHA: change.Before, Hard: true})
	}

	// re-create locally removed branches
	for _, removedLocalBranch := range self.LocalRemoved.BranchNames() {
		startingPoint := self.LocalRemoved[removedLocalBranch]
		result.Add(&opcode.CreateBranch{
			Branch:        removedLocalBranch,
			StartingPoint: startingPoint.Location(),
		})
	}

	// remove locally added branches
	for _, addedLocalBranch := range self.LocalAdded {
		if args.FinalBranch == addedLocalBranch {
			result.Add(&opcode.Checkout{Branch: args.InitialBranch})
		}
		result.Add(&opcode.DeleteLocalBranch{
			Branch: addedLocalBranch,
			Force:  true,
		})
	}

	// Ignore remotely changed perennial branches because we can't force-push to them
	// and we would need the local branch to revert commits on them, but we can't change the local branch.

	// reset remotely changed feature branches
	_, remoteFeatureChanges := CategorizeRemoteBranchChange(self.RemoteChanged, args.BranchTypes)
	for _, remoteChangedFeatureBranch := range remoteFeatureChanges.BranchNames() {
		change := remoteFeatureChanges[remoteChangedFeatureBranch]
		result.Add(&opcode.ResetRemoteBranchToSHA{
			Branch:      remoteChangedFeatureBranch,
			MustHaveSHA: change.After,
			SetToSHA:    change.Before,
		})
	}

	// This must be a CheckoutIfExists opcode because this branch might not exist
	// when a Git Town command fails, stores this undo opcode, then gets continued and deletes this branch.
	result.Add(&opcode.CheckoutIfExists{Branch: args.InitialBranch})
	return result
}

type BranchChangesUndoProgramArgs struct {
	Lineage                  configdomain.Lineage
	BranchTypes              configdomain.BranchTypes
	InitialBranch            gitdomain.LocalBranchName
	FinalBranch              gitdomain.LocalBranchName
	NoPushHook               configdomain.NoPushHook
	UndoablePerennialCommits []gitdomain.SHA
}
