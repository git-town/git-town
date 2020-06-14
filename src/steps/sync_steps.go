package steps

import (
	"fmt"

	"github.com/git-town/git-town/src/git"
)

// GetSyncBranchSteps returns the steps to sync the branch with the given name.
func GetSyncBranchSteps(branchName string, pushBranch bool, repo *git.ProdRepo) (result StepList, err error) {
	isFeature := repo.Config.IsFeatureBranch(branchName)
	hasRemoteOrigin, err := repo.Silent.HasRemote("origin")
	if err != nil {
		return result, err
	}
	if !hasRemoteOrigin && !isFeature {
		return
	}
	result.Append(&CheckoutBranchStep{BranchName: branchName})
	if isFeature {
		steps, err := getSyncFeatureBranchSteps(branchName, repo)
		if err != nil {
			return result, err
		}
		result.AppendList(steps)
	} else {
		steps, err := getSyncNonFeatureBranchSteps(branchName, repo)
		if err != nil {
			return result, err
		}
		result.AppendList(steps)
	}
	if pushBranch && hasRemoteOrigin && !repo.Config.IsOffline() {
		hasTrackingBranch, err := repo.Silent.HasTrackingBranch(branchName)
		if err != nil {
			return result, err
		}
		if hasTrackingBranch {
			result.Append(&PushBranchStep{BranchName: branchName})
		} else {
			result.Append(&CreateTrackingBranchStep{BranchName: branchName})
		}
	}
	return result, nil
}

// Helpers

func getSyncFeatureBranchSteps(branchName string, repo *git.ProdRepo) (result StepList, err error) {
	hasTrackingBranch, err := repo.Silent.HasTrackingBranch(branchName)
	if err != nil {
		return result, err
	}
	if hasTrackingBranch {
		result.Append(&MergeBranchStep{BranchName: repo.Silent.TrackingBranchName(branchName)})
	}
	result.Append(&MergeBranchStep{BranchName: repo.Config.GetParentBranch(branchName)})
	return
}

func getSyncNonFeatureBranchSteps(branchName string, repo *git.ProdRepo) (result StepList, err error) {
	hasTrackingBranch, err := repo.Silent.HasTrackingBranch(branchName)
	if err != nil {
		return result, err
	}
	if hasTrackingBranch {
		if repo.Config.GetPullBranchStrategy() == "rebase" {
			result.Append(&RebaseBranchStep{BranchName: repo.Silent.TrackingBranchName(branchName)})
		} else {
			result.Append(&MergeBranchStep{BranchName: repo.Silent.TrackingBranchName(branchName)})
		}
	}

	mainBranchName := repo.Config.GetMainBranch()
	hasUpstream, err := repo.Silent.HasRemote("upstream")
	if err != nil {
		return result, err
	}
	if mainBranchName == branchName && hasUpstream && repo.Config.ShouldSyncUpstream() {
		result.Append(&FetchUpstreamStep{BranchName: mainBranchName})
		result.Append(&RebaseBranchStep{BranchName: fmt.Sprintf("upstream/%s", mainBranchName)})
	}
	return
}
