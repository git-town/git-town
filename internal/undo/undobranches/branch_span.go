package undobranches

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
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
		if localName, hasLocalName := before.LocalName.Get(); hasLocalName {
			branchNames.Add(localName.BranchName())
		}
		if remoteName, hasRemoteName := before.RemoteName.Get(); hasRemoteName {
			branchNames.Add(remoteName.BranchName())
		}
	}
	if after, hasAfter := self.After.Get(); hasAfter {
		if localName, hasLocalName := after.LocalName.Get(); hasLocalName {
			branchNames.Add(localName.BranchName())
		}
		if remoteName, hasRemoteName := after.RemoteName.Get(); hasRemoteName {
			branchNames.Add(remoteName.BranchName())
		}
	}
	return branchNames.Values()
}

func (self BranchSpan) IsInconsistentChange() Option[InconsistentChange] {
	omniChangeData := self.IsOmniChange()
	localChangedResult := self.LocalChanged()
	remoteChanged := self.RemoteChanged()
	before, hasBefore := self.Before.Get()
	after, hasAfter := self.After.Get()
	isInconsistentChange := hasBefore && before.HasTrackingBranch() && hasAfter && after.HasTrackingBranch() && localChangedResult.IsChanged && remoteChanged.IsChanged && !omniChangeData.IsOmniChange
	if isInconsistentChange {
		return Some(InconsistentChange{
			Before: before,
			After:  after,
		})
	} else {
		return None[InconsistentChange]()
	}
}

type InconsistentChange struct {
	Before gitdomain.BranchInfo
	After  gitdomain.BranchInfo
}

// IsLocalRename indicates whether this BranchSpan describes the situation where only the local branch was renamed.
func (self BranchSpan) IsLocalRename() (isLocalRename bool, beforeName, afterName gitdomain.LocalBranchName) {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return false, "", ""
	}
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return false, "", ""
	}
	beforeName, hasBeforeName := before.LocalName.Get()
	if !hasBeforeName {
		return false, "", ""
	}
	afterName, hasAfterName := after.LocalName.Get()
	if !hasAfterName {
		return false, "", ""
	}
	beforeSHA, hasBeforeSHA := before.LocalSHA.Get()
	if !hasBeforeSHA {
		return false, "", ""
	}
	afterSHA, hasAfterSHA := after.LocalSHA.Get()
	if !hasAfterSHA {
		return false, "", ""
	}
	return beforeName != afterName && beforeSHA == afterSHA, beforeName, afterName
}

// IsOmniChange indicates whether this BranchBeforeAfter changes a synced branch
// from one SHA both locally and remotely to another SHA both locally and remotely.
func (self BranchSpan) IsOmniChange() IsOmniChangeResult {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return IsOmniChangeResult{
			IsOmniChange: false,
			Name:         "",
			SHAAfter:     "",
			SHABefore:    "",
		}
	}
	beforeIsOmni, beforeName, beforeSHA := before.IsOmniBranch()
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return IsOmniChangeResult{
			IsOmniChange: false,
			Name:         "",
			SHAAfter:     "",
			SHABefore:    "",
		}
	}
	afterIsOmni, _, afterSHA := after.IsOmniBranch()
	return IsOmniChangeResult{
		IsOmniChange: beforeIsOmni && afterIsOmni && beforeSHA != afterSHA,
		Name:         beforeName,
		SHAAfter:     afterSHA,
		SHABefore:    beforeSHA,
	}
}

type IsOmniChangeResult struct {
	IsOmniChange bool
	Name         gitdomain.LocalBranchName
	SHAAfter     gitdomain.SHA
	SHABefore    gitdomain.SHA
}

// Indicates whether this BranchSpan describes the removal of an omni Branch
// and provides all relevant data for this situation.
func (self BranchSpan) IsOmniRemove() IsOmniRemoveResult {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return IsOmniRemoveResult{
			IsOmniRemove: false,
			Name:         "",
			SHA:          "",
		}
	}
	beforeIsOmni, beforeName, beforeSHA := before.IsOmniBranch()
	_, hasAfter := self.After.Get()
	return IsOmniRemoveResult{
		IsOmniRemove: beforeIsOmni && !hasAfter,
		Name:         beforeName,
		SHA:          beforeSHA,
	}
}

type IsOmniRemoveResult struct {
	IsOmniRemove bool
	Name         gitdomain.LocalBranchName
	SHA          gitdomain.SHA
}

func (self BranchSpan) LocalAdded() LocalAddedResult {
	before, hasBefore := self.Before.Get()
	beforeHasLocalBranch, _, _ := before.GetLocal()
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return LocalAddedResult{
			IsAdded: false,
			Name:    "",
			SHA:     "",
		}
	}
	afterHasLocalBranch, afterLocalBranch, afterSHA := after.GetLocal()
	return LocalAddedResult{
		IsAdded: (!hasBefore || !beforeHasLocalBranch) && afterHasLocalBranch,
		Name:    afterLocalBranch,
		SHA:     afterSHA,
	}
}

type LocalAddedResult struct {
	IsAdded bool
	Name    gitdomain.LocalBranchName
	SHA     gitdomain.SHA
}

func (self BranchSpan) LocalChanged() LocalChangedResult {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return LocalChangedResult{
			IsChanged: false,
			Name:      "",
			SHAAfter:  "",
			SHABefore: "",
		}
	}
	hasLocalBranchBefore, beforeBranch, beforeSHA := before.GetLocal()
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return LocalChangedResult{
			IsChanged: false,
			Name:      "",
			SHAAfter:  "",
			SHABefore: "",
		}
	}
	hasLocalBranchAfter, _, afterSHA := after.GetLocal()
	return LocalChangedResult{
		IsChanged: hasLocalBranchBefore && hasLocalBranchAfter && beforeSHA != afterSHA,
		Name:      beforeBranch,
		SHAAfter:  afterSHA,
		SHABefore: beforeSHA,
	}
}

type LocalChangedResult struct {
	IsChanged bool
	Name      gitdomain.LocalBranchName
	SHAAfter  gitdomain.SHA
	SHABefore gitdomain.SHA
}

func (self BranchSpan) LocalRemoved() LocalRemovedResult {
	before, hasBefore := self.Before.Get()
	hasBeforeBranch, branchName, beforeSHA := before.GetLocal()
	after, hasAfter := self.After.Get()
	hasAfterBranch, _, _ := after.GetLocal()
	return LocalRemovedResult{
		IsRemoved: hasBefore && hasBeforeBranch && (!hasAfter || !hasAfterBranch),
		Name:      branchName,
		SHA:       beforeSHA,
	}
}

type LocalRemovedResult struct {
	IsRemoved bool
	Name      gitdomain.LocalBranchName
	SHA       gitdomain.SHA
}

func (self BranchSpan) RemoteAdded() RemoteAddedResult {
	before, hasBefore := self.Before.Get()
	beforeHasRemoteBranch, _, _ := before.GetRemote()
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return RemoteAddedResult{
			IsAdded: false,
			Name:    "",
			SHA:     "",
		}
	}
	afterHasRemoteBranch, afterRemoteBranchName, afterRemoteBranchSHA := after.GetRemote()
	return RemoteAddedResult{
		IsAdded: (!hasBefore || !beforeHasRemoteBranch) && afterHasRemoteBranch,
		Name:    afterRemoteBranchName,
		SHA:     afterRemoteBranchSHA,
	}
}

type RemoteAddedResult struct {
	IsAdded bool
	Name    gitdomain.RemoteBranchName
	SHA     gitdomain.SHA
}

func (self BranchSpan) RemoteChanged() RemoteChangedResult {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return RemoteChangedResult{
			Name:      "",
			IsChanged: false,
			SHAAfter:  "",
			SHABefore: "",
		}
	}
	beforeHasRemoteBranch, beforeRemoteBranchName, beforeRemoteBranchSHA := before.GetRemote()
	after, hasAfter := self.After.Get()
	if !hasAfter {
		return RemoteChangedResult{
			Name:      "",
			IsChanged: false,
			SHAAfter:  "",
			SHABefore: "",
		}
	}
	afterHasRemoteBranch, _, afterRemoteBranchSHA := after.GetRemote()
	return RemoteChangedResult{
		Name:      beforeRemoteBranchName,
		IsChanged: beforeHasRemoteBranch && afterHasRemoteBranch && beforeRemoteBranchSHA != afterRemoteBranchSHA,
		SHAAfter:  afterRemoteBranchSHA,
		SHABefore: beforeRemoteBranchSHA,
	}
}

type RemoteChangedResult struct {
	Name      gitdomain.RemoteBranchName
	IsChanged bool
	SHAAfter  gitdomain.SHA
	SHABefore gitdomain.SHA
}

func (self BranchSpan) RemoteRemoved() RemoteRemovedResult {
	before, hasBefore := self.Before.Get()
	if !hasBefore {
		return RemoteRemovedResult{
			IsRemoved: false,
			Name:      "",
			SHA:       "",
		}
	}
	beforeHasRemoteBranch, remoteBranchName, beforeSHA := before.GetRemote()
	after, hasAfter := self.After.Get()
	afterHasRemoteBranch, _, _ := after.GetRemote()
	return RemoteRemovedResult{
		IsRemoved: beforeHasRemoteBranch && (!hasAfter || !afterHasRemoteBranch),
		Name:      remoteBranchName,
		SHA:       beforeSHA,
	}
}

type RemoteRemovedResult struct {
	IsRemoved bool
	Name      gitdomain.RemoteBranchName
	SHA       gitdomain.SHA
}

// func (self BranchSpan) String() string {
// 	result := "BranchSpan:\n"
// 	result += "Before:" + self.Before.String() + "\n"
// 	result += "After:" + self.After.String() + "\n"
// 	return result
// }
