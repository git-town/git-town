package steps

import (
	"fmt"

	"github.com/Originate/git-town/src/git"
)

// GetSyncBranchSteps returns the steps to sync the branch with the given name.
func GetSyncBranchSteps(branchName string) (result StepList) {
	isFeature := git.IsFeatureBranch(branchName)
	hasRemoteOrigin := git.HasRemote("origin")

	if !hasRemoteOrigin && !isFeature {
		return
	}

	result.Append(&CheckoutBranchStep{BranchName: branchName})
	if isFeature {
		if git.HasTrackingBranch(branchName) {
			result.Append(&MergeBranchStep{BranchName: git.GetTrackingBranchName(branchName)})
		}
		result.Append(&MergeBranchStep{BranchName: git.GetParentBranch(branchName)})
	} else {
		if git.HasTrackingBranch(branchName) {
			if git.GetPullBranchStrategy() == "rebase" {
				result.Append(&RebaseBranchStep{BranchName: git.GetTrackingBranchName(branchName)})
			} else {
				result.Append(&MergeBranchStep{BranchName: git.GetTrackingBranchName(branchName)})
			}
		}

		mainBranchName := git.GetMainBranch()
		if mainBranchName == branchName && git.HasRemote("upstream") {
			result.Append(&FetchUpstreamStep{})
			result.Append(&RebaseBranchStep{BranchName: fmt.Sprintf("upstream/%s", mainBranchName)})
		}
	}

	if hasRemoteOrigin && !git.IsOffline() {
		if git.HasTrackingBranch(branchName) {
			result.Append(&PushBranchStep{BranchName: branchName})
		} else {
			result.Append(&CreateTrackingBranchStep{BranchName: branchName})
		}
	}

	return
}
