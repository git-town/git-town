package config

import (
	"errors"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// BranchesToMark provides the branches to make contribution, observed, parked, or prototype.
func BranchesToMark(args []string, branchesSnapshot gitdomain.BranchesSnapshot, config UnvalidatedConfig) (BranchesToMarkResult, error) {
	branchesToMark := configdomain.BranchesAndTypes{}
	var branchToCheckout Option[gitdomain.LocalBranchName]
	var emptyResult BranchesToMarkResult
	switch len(args) {
	case 0:
		currentBranch, hasCurrentBranch := branchesSnapshot.Active.Get()
		if !hasCurrentBranch {
			return emptyResult, errors.New(messages.CurrentBranchCannotDetermine)
		}
		branchesToMark.AddTypeFor(currentBranch, &config)
		branchToCheckout = None[gitdomain.LocalBranchName]()
	case 1:
		branch := gitdomain.NewLocalBranchName(args[0])
		branchesToMark.AddTypeFor(branch, &config)
		branchInfo, hasBranchInfo := branchesSnapshot.Branches.FindRemoteNameMatchingLocal(branch).Get()
		if hasBranchInfo && branchInfo.SyncStatus == gitdomain.SyncStatusRemoteOnly {
			branchToCheckout = Some(branch)
		} else {
			branchToCheckout = None[gitdomain.LocalBranchName]()
		}
	default:
		branchesToMark.AddMany(gitdomain.NewLocalBranchNames(args...), &config)
		branchToCheckout = None[gitdomain.LocalBranchName]()
	}
	return BranchesToMarkResult{
		BranchToCheckout: branchToCheckout,
		BranchesToMark:   branchesToMark,
	}, nil
}

type BranchesToMarkResult struct {
	BranchToCheckout Option[gitdomain.LocalBranchName]
	BranchesToMark   configdomain.BranchesAndTypes
}
