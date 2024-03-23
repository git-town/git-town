package undobranches

import (
	"github.com/git-town/git-town/v13/src/config/configdomain"
	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/src/gohacks/slice"
	"github.com/git-town/git-town/v13/src/undo/undodomain"
	"github.com/git-town/git-town/v13/src/vm/opcodes"
	"github.com/git-town/git-town/v13/src/vm/program"
)

// BranchChanges describes the changes made to the branches in a Git repo.
// Various types of changes are distinguished.
type BranchChanges struct {
	// Inconsistent changes are changes on both local and tracking branch, but where the local and tracking branch
	// don't have the same SHA before or after.
	// These changes cannot be undone for perennial branches because there is no way to reset the remote branch to the SHA it had before.
	InconsistentlyChanged undodomain.InconsistentChanges
	LocalAdded            gitdomain.LocalBranchNames
	LocalChanged          LocalBranchChange
	LocalRemoved          LocalBranchesSHAs
	// OmniChanges are changes where the local SHA and the remote SHA are identical before the change as well as after the change,
	OmniChanged LocalBranchChange // a branch had the same SHA locally and remotely, now it has a new SHA locally and remotely, the local and remote SHA are still equal
	// OmniRemoved is when a branch that has the same SHA on its local and tracking branch gets removed.
	OmniRemoved   LocalBranchesSHAs
	RemoteAdded   gitdomain.RemoteBranchNames
	RemoteChanged RemoteBranchChange
	RemoteRemoved RemoteBranchesSHAs
}

// EmptyBranchChanges provides a properly initialized empty Changes instance.
func EmptyBranchChanges() BranchChanges {
	return BranchChanges{
		InconsistentlyChanged: undodomain.InconsistentChanges{},
		LocalAdded:            gitdomain.LocalBranchNames{},
		LocalChanged:          LocalBranchChange{},
		LocalRemoved:          LocalBranchesSHAs{},
		OmniChanged:           LocalBranchChange{},
		OmniRemoved:           LocalBranchesSHAs{},
		RemoteAdded:           gitdomain.RemoteBranchNames{},
		RemoteChanged:         map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]{},
		RemoteRemoved:         map[gitdomain.RemoteBranchName]gitdomain.SHA{},
	}
}

// uncomment when debugging the undo logic
// func (self BranchChanges) String() string {
// 	s := strings.Builder{}
// 	s.WriteString("BranchChanges {")
// 	s.WriteString("\n  LocalAdded: ")
// 	s.WriteString(strings.Join(self.LocalAdded.Strings(), ", "))
// 	s.WriteString("\n  LocalRemoved: ")
// 	s.WriteString(strings.Join(self.LocalRemoved.BranchNames().Strings(), ", "))
// 	s.WriteString("\n  LocalChanged: ")
// 	s.WriteString(strings.Join(self.LocalChanged.BranchNames().Strings(), ", "))
// 	s.WriteString("\n  RemoteAdded: ")
// 	s.WriteString(strings.Join(self.RemoteAdded.Strings(), ", "))
// 	s.WriteString("\n  RemoteRemoved: ")
// 	s.WriteString(strings.Join(self.RemoteRemoved.BranchNames().Strings(), ", "))
// 	s.WriteString("\n  RemoteChanged: ")
// 	s.WriteString(strings.Join(self.RemoteChanged.BranchNames().Strings(), ", "))
// 	s.WriteString("\n  OmniRemoved: ")
// 	s.WriteString(strings.Join(self.OmniRemoved.BranchNames().Strings(), ", "))
// 	s.WriteString("\n  OmniChanged: ")
// 	s.WriteString(strings.Join(self.OmniChanged.BranchNames().Strings(), ", "))
// 	s.WriteString("\n  InconsistentlyChanged: ")
// 	s.WriteString(strings.Join(self.InconsistentlyChanged.BranchNames().Strings(), ", "))
// 	s.WriteRune('\n')
// 	return s.String()
// }

// UndoProgram provides the steps to undo the changes described by this BranchChanges instance.
func (self BranchChanges) UndoProgram(args BranchChangesUndoProgramArgs) program.Program {
	result := program.Program{}
	omniChangedPerennials, omniChangedFeatures := CategorizeLocalBranchChange(self.OmniChanged, args.Config)

	// revert omni-changed perennial branches
	for _, branch := range omniChangedPerennials.BranchNames() {
		change := omniChangedPerennials[branch]
		if slice.Contains(args.UndoablePerennialCommits, change.After) {
			result.Add(&opcodes.Checkout{Branch: branch})
			result.Add(&opcodes.RevertCommit{SHA: change.After})
			result.Add(&opcodes.PushCurrentBranch{CurrentBranch: branch})
		}
	}

	// reset omni-changed feature branches
	for _, branch := range omniChangedFeatures.BranchNames() {
		change := omniChangedFeatures[branch]
		result.Add(&opcodes.Checkout{Branch: branch})
		result.Add(&opcodes.ResetCurrentBranchToSHA{MustHaveSHA: change.After, SetToSHA: change.Before, Hard: true})
		result.Add(&opcodes.ForcePushCurrentBranch{})
	}

	// re-create removed omni-branches
	for _, branch := range self.OmniRemoved.BranchNames() {
		sha := self.OmniRemoved[branch]
		result.Add(&opcodes.CreateBranch{Branch: branch, StartingPoint: sha.Location()})
		result.Add(&opcodes.CreateTrackingBranch{Branch: branch})
	}

	inconsistentlyChangedPerennials, inconsistentChangedFeatures := CategorizeInconsistentChanges(self.InconsistentlyChanged, args.Config)

	// reset inconsintently changed perennial branches
	for _, inconsistentlyChangedPerennial := range inconsistentlyChangedPerennials {
		if inconsistentlyChangedPerennial.After.IsOmniBranch() {
			if slice.Contains(args.UndoablePerennialCommits, inconsistentlyChangedPerennial.After.LocalSHA) {
				result.Add(&opcodes.Checkout{Branch: inconsistentlyChangedPerennial.Before.LocalName})
				result.Add(&opcodes.RevertCommit{SHA: inconsistentlyChangedPerennial.After.LocalSHA})
				result.Add(&opcodes.PushCurrentBranch{CurrentBranch: inconsistentlyChangedPerennial.After.LocalName})
			}
		}
	}

	// reset inconsintently changed feature branches
	for _, inconsistentChange := range inconsistentChangedFeatures {
		result.Add(&opcodes.Checkout{Branch: inconsistentChange.Before.LocalName})
		result.Add(&opcodes.ResetCurrentBranchToSHA{
			MustHaveSHA: inconsistentChange.After.LocalSHA,
			SetToSHA:    inconsistentChange.Before.LocalSHA,
			Hard:        true,
		})
		result.Add(&opcodes.ResetRemoteBranchToSHA{
			Branch:      inconsistentChange.Before.RemoteName,
			MustHaveSHA: inconsistentChange.After.RemoteSHA,
			SetToSHA:    inconsistentChange.Before.RemoteSHA,
		})
	}

	// remove remotely added branches
	for _, addedRemoteBranch := range self.RemoteAdded {
		if addedRemoteBranch.Remote() != gitdomain.UpstreamRemote {
			result.Add(&opcodes.DeleteTrackingBranch{
				Branch: addedRemoteBranch,
			})
		}
	}

	// re-create remotely removed feature branches
	_, removedFeatureTrackingBranches := CategorizeRemoteBranchesSHAs(self.RemoteRemoved, args.Config)
	for _, branch := range removedFeatureTrackingBranches.BranchNames() {
		sha := removedFeatureTrackingBranches[branch]
		result.Add(&opcodes.CreateRemoteBranch{
			Branch: branch.LocalBranchName(),
			SHA:    sha,
		})
	}

	// reset locally changed branches
	for _, localBranch := range self.LocalChanged.BranchNames() {
		change := self.LocalChanged[localBranch]
		result.Add(&opcodes.Checkout{Branch: localBranch})
		result.Add(&opcodes.ResetCurrentBranchToSHA{MustHaveSHA: change.After, SetToSHA: change.Before, Hard: true})
	}

	// re-create locally removed branches
	for _, removedLocalBranch := range self.LocalRemoved.BranchNames() {
		startingPoint := self.LocalRemoved[removedLocalBranch]
		result.Add(&opcodes.CreateBranch{
			Branch:        removedLocalBranch,
			StartingPoint: startingPoint.Location(),
		})
	}

	// remove locally added branches
	for _, addedLocalBranch := range self.LocalAdded {
		if args.EndBranch == addedLocalBranch {
			result.Add(&opcodes.Checkout{Branch: args.BeginBranch})
		}
		result.Add(&opcodes.DeleteLocalBranch{Branch: addedLocalBranch})
	}

	// Ignore remotely changed perennial branches because we can't force-push to them
	// and we would need the local branch to revert commits on them, but we can't change the local branch.

	// reset remotely changed feature branches
	_, remoteFeatureChanges := CategorizeRemoteBranchChange(self.RemoteChanged, args.Config)
	for _, remoteChangedFeatureBranch := range remoteFeatureChanges.BranchNames() {
		change := remoteFeatureChanges[remoteChangedFeatureBranch]
		result.Add(&opcodes.ResetRemoteBranchToSHA{
			Branch:      remoteChangedFeatureBranch,
			MustHaveSHA: change.After,
			SetToSHA:    change.Before,
		})
	}

	// This must be a CheckoutIfExists opcode because this branch might not exist
	// when a Git Town command fails, stores this undo opcode, then gets continued and deletes this branch.
	result.Add(&opcodes.CheckoutIfExists{Branch: args.BeginBranch})
	return result
}

type BranchChangesUndoProgramArgs struct {
	BeginBranch              gitdomain.LocalBranchName
	Config                   *configdomain.FullConfig
	EndBranch                gitdomain.LocalBranchName
	UndoablePerennialCommits []gitdomain.SHA
}
