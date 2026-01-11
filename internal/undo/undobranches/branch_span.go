package undobranches

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/undo/undodomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/git-town/git-town/v22/pkg/set"
)

// BranchSpan represents changes of a branch over time.
type BranchSpan struct {
	Before Option[gitdomain.BranchInfo] // the status of the branch before Git Town ran
	After  Option[gitdomain.BranchInfo] // the status of the branch after Git Town ran
}

func (self BranchSpan) BranchNames() []gitdomain.BranchName {
	branchNames := set.New[gitdomain.BranchName]()
	if before, hasBefore := self.Before.Get(); hasBefore {
		if localName, hasLocalName := before.LocalName().Get(); hasLocalName {
			branchNames.Add(localName.BranchName())
		}
		if remoteName, hasRemoteName := before.RemoteName.Get(); hasRemoteName {
			branchNames.Add(remoteName.BranchName())
		}
	}
	if after, hasAfter := self.After.Get(); hasAfter {
		if localName, hasLocalName := after.LocalName().Get(); hasLocalName {
			branchNames.Add(localName.BranchName())
		}
		if remoteName, hasRemoteName := after.RemoteName.Get(); hasRemoteName {
			branchNames.Add(remoteName.BranchName())
		}
	}
	return branchNames.Values()
}

func (self BranchSpan) InconsistentChange() Option[undodomain.InconsistentChange] {
	_, isOmniChange := self.OmniChange().Get()
	_, localChanged := self.LocalChange().Get()
	_, remoteChanged := self.RemoteChange().Get()
	before, hasBefore := self.Before.Get()
	after, hasAfter := self.After.Get()
	isInconsistentChange := hasBefore && before.HasTrackingBranch() && hasAfter && after.HasTrackingBranch() && localChanged && remoteChanged && !isOmniChange
	if !isInconsistentChange {
		return None[undodomain.InconsistentChange]()
	}
	return Some(undodomain.InconsistentChange{
		After:  after,
		Before: before,
	})
}

func (self BranchSpan) LocalAdd() Option[gitdomain.LocalBranchName] {
	before, hasBefore := self.Before.Get()
	_, beforeHasLocalBranch := before.Local.Get()
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return None[gitdomain.LocalBranchName]()
	}
	afterLocalBranch, afterHasLocalBranch := after.Local.Get()
	isLocalAdded := (!hasBefore || !beforeHasLocalBranch) && afterHasLocalBranch
	if !isLocalAdded {
		return None[gitdomain.LocalBranchName]()
	}
	return Some(afterLocalBranch.Name)
}

func (self BranchSpan) LocalChange() Option[LocalBranchChange] {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return None[LocalBranchChange]()
	}
	beforeLocal, hasLocalBefore := before.Local.Get()
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return None[LocalBranchChange]()
	}
	afterLocal, hasLocalBranchAfter := after.Local.Get()
	localChanged := hasLocalBefore && hasLocalBranchAfter && beforeLocal.SHA != afterLocal.SHA
	if !localChanged {
		return None[LocalBranchChange]()
	}
	return Some(LocalBranchChange{
		beforeLocal.Name: undodomain.Change[gitdomain.SHA]{
			Before: beforeLocal.SHA,
			After:  afterLocal.SHA,
		},
	})
}

func (self BranchSpan) LocalRemove() Option[LocalBranchesSHAs] {
	before, hasBefore := self.Before.Get()
	beforeLocal, hasBeforeLocal := before.Local.Get()
	after, hasAfter := self.After.Get()
	_, hasAfterBranch := after.Local.Get()
	localRemoved := hasBefore && hasBeforeLocal && (!hasAfter || !hasAfterBranch)
	if !localRemoved {
		return None[LocalBranchesSHAs]()
	}
	return Some(LocalBranchesSHAs{
		beforeLocal.Name: beforeLocal.SHA,
	})
}

// LocalRename indicates whether this BranchSpan describes the situation where only the local branch was renamed.
func (self BranchSpan) LocalRename() Option[LocalBranchRename] {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return None[LocalBranchRename]()
	}
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return None[LocalBranchRename]()
	}
	beforeName, hasBeforeName := before.LocalName.Get()
	if !hasBeforeName {
		return None[LocalBranchRename]()
	}
	afterName, hasAfterName := after.LocalName.Get()
	if !hasAfterName {
		return None[LocalBranchRename]()
	}
	beforeSHA, hasBeforeSHA := before.LocalSHA.Get()
	if !hasBeforeSHA {
		return None[LocalBranchRename]()
	}
	afterSHA, hasAfterSHA := after.LocalSHA.Get()
	if !hasAfterSHA {
		return None[LocalBranchRename]()
	}
	isLocalRename := beforeName != afterName && beforeSHA == afterSHA
	if !isLocalRename {
		return None[LocalBranchRename]()
	}
	return Some(LocalBranchRename{
		After:  afterName,
		Before: beforeName,
	})
}

// OmniChange indicates whether this BranchBeforeAfter changes a synced branch
// from one SHA both locally and remotely to another SHA both locally and remotely.
func (self BranchSpan) OmniChange() Option[LocalBranchChange] {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return None[LocalBranchChange]()
	}
	beforeIsOmni, beforeName, beforeSHA := before.IsOmniBranch()
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return None[LocalBranchChange]()
	}
	afterIsOmni, _, afterSHA := after.IsOmniBranch()
	isOmniChange := beforeIsOmni && afterIsOmni && beforeSHA != afterSHA
	if !isOmniChange {
		return None[LocalBranchChange]()
	}
	return Some(LocalBranchChange{
		beforeName: undodomain.Change[gitdomain.SHA]{
			Before: beforeSHA,
			After:  afterSHA,
		},
	})
}

// Indicates whether this BranchSpan describes the removal of an omni Branch
// and provides all relevant data for this situation.
func (self BranchSpan) OmniRemove() Option[LocalBranchesSHAs] {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return None[LocalBranchesSHAs]()
	}
	beforeIsOmni, beforeName, beforeSHA := before.IsOmniBranch()
	_, hasAfter := self.After.Get()
	isOmniRemove := beforeIsOmni && !hasAfter
	if !isOmniRemove {
		return None[LocalBranchesSHAs]()
	}
	return Some(LocalBranchesSHAs{
		beforeName: beforeSHA,
	})
}

func (self BranchSpan) RemoteAdd() Option[gitdomain.RemoteBranchName] {
	before, hasBefore := self.Before.Get()
	beforeHasRemoteBranch, _, _ := before.GetRemote()
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return None[gitdomain.RemoteBranchName]()
	}
	afterHasRemoteBranch, afterRemoteBranchName, _ := after.GetRemote()
	remoteAdded := (!hasBefore || !beforeHasRemoteBranch) && afterHasRemoteBranch
	if !remoteAdded {
		return None[gitdomain.RemoteBranchName]()
	}
	return Some(afterRemoteBranchName)
}

func (self BranchSpan) RemoteChange() Option[RemoteBranchChange] {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return None[RemoteBranchChange]()
	}
	beforeHasRemoteBranch, beforeRemoteBranchName, beforeRemoteBranchSHA := before.GetRemote()
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return None[RemoteBranchChange]()
	}
	afterHasRemoteBranch, _, afterRemoteBranchSHA := after.GetRemote()
	remoteChanged := beforeHasRemoteBranch && afterHasRemoteBranch && beforeRemoteBranchSHA != afterRemoteBranchSHA
	if !remoteChanged {
		return None[RemoteBranchChange]()
	}
	return Some(RemoteBranchChange{
		beforeRemoteBranchName: undodomain.Change[gitdomain.SHA]{
			Before: beforeRemoteBranchSHA,
			After:  afterRemoteBranchSHA,
		},
	})
}

func (self BranchSpan) RemoteRemove() Option[RemoteBranchesSHAs] {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return None[RemoteBranchesSHAs]()
	}
	beforeHasRemoteBranch, remoteBranchName, beforeSHA := before.GetRemote()
	after, hasAfter := self.After.Get()
	afterHasRemoteBranch, _, _ := after.GetRemote()
	remoteRemoved := beforeHasRemoteBranch && (!hasAfter || !afterHasRemoteBranch)
	if !remoteRemoved {
		return None[RemoteBranchesSHAs]()
	}
	return Some(RemoteBranchesSHAs{
		remoteBranchName: beforeSHA,
	})
}

// func (self BranchSpan) String() string {
// 	result := "BranchSpan:\n"
// 	result += "Before:" + self.Before.String() + "\n"
// 	result += "After:" + self.After.String() + "\n"
// 	return result
// }
