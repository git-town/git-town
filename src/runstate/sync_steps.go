package runstate

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/steps"
)

// SyncBranchSteps provides the steps to sync the branch with the given name.
func SyncBranchSteps(branchName string, pushBranch bool, repo *git.ProdRepo) (result StepList, err error) {
	isFeature := repo.Config.IsFeatureBranch(branchName)
	hasOrigin, err := repo.Silent.HasOrigin()
	if err != nil {
		return result, err
	}
	if !hasOrigin && !isFeature {
		return
	}
	result.Append(&steps.CheckoutBranchStep{BranchName: branchName})
	if isFeature {
		steps, err := syncFeatureBranchSteps(branchName, repo)
		if err != nil {
			return result, err
		}
		result.AppendList(steps)
	} else {
		steps, err := syncNonFeatureBranchSteps(branchName, repo)
		if err != nil {
			return result, err
		}
		result.AppendList(steps)
	}
	if pushBranch && hasOrigin && !repo.Config.IsOffline() {
		hasTrackingBranch, err := repo.Silent.HasTrackingBranch(branchName)
		if err != nil {
			return result, err
		}
		if hasTrackingBranch {
			result.Append(&steps.PushBranchStep{BranchName: branchName})
		} else {
			result.Append(&steps.CreateTrackingBranchStep{BranchName: branchName})
		}
	}
	return result, nil
}

// Helpers

func syncFeatureBranchSteps(branchName string, repo *git.ProdRepo) (StepList, error) {
	hasTrackingBranch, err := repo.Silent.HasTrackingBranch(branchName)
	if err != nil {
		return StepList{}, err
	}
	result := StepList{}
	if hasTrackingBranch {
		result.Append(&steps.MergeBranchStep{BranchName: repo.Silent.TrackingBranchName(branchName)})
	}
	result.Append(&steps.MergeBranchStep{BranchName: repo.Config.ParentBranch(branchName)})
	return result, nil
}

func syncNonFeatureBranchSteps(branchName string, repo *git.ProdRepo) (StepList, error) {
	hasTrackingBranch, err := repo.Silent.HasTrackingBranch(branchName)
	if err != nil {
		return StepList{}, err
	}
	result := StepList{}
	if hasTrackingBranch {
		if repo.Config.PullBranchStrategy() == "rebase" {
			result.Append(&steps.RebaseBranchStep{BranchName: repo.Silent.TrackingBranchName(branchName)})
		} else {
			result.Append(&steps.MergeBranchStep{BranchName: repo.Silent.TrackingBranchName(branchName)})
		}
	}

	mainBranchName := repo.Config.MainBranch()
	hasUpstream, err := repo.Silent.HasRemote("upstream")
	if err != nil {
		return result, err
	}
	if mainBranchName == branchName && hasUpstream && repo.Config.ShouldSyncUpstream() {
		result.Append(&steps.FetchUpstreamStep{BranchName: mainBranchName})
		result.Append(&steps.RebaseBranchStep{BranchName: fmt.Sprintf("upstream/%s", mainBranchName)})
	}
	return result, nil
}
