package steps

import (
	"fmt"

	"github.com/git-town/git-town/src/git"
)

// GetSyncBranchSteps returns the steps to sync the branch with the given name.
func GetSyncBranchSteps(branchName string, pushBranch bool) (result StepList) {
	isFeature := git.Config().IsFeatureBranch(branchName)
	hasRemoteOrigin := git.HasRemote("origin")

	if !hasRemoteOrigin && !isFeature {
		return
	}

	result.Append(&CheckoutBranchStep{BranchName: branchName})
	if isFeature {
		result.AppendList(getSyncFeatureBranchSteps(branchName))
	} else {
		result.AppendList(getSyncNonFeatureBranchSteps(branchName))
	}

	if pushBranch && hasRemoteOrigin && !git.Config().IsOffline() {
		if git.HasTrackingBranch(branchName) {
			result.Append(&PushBranchStep{BranchName: branchName})
		} else {
			result.Append(&CreateTrackingBranchStep{BranchName: branchName})
		}
	}

	return
}

// Helpers

func getSyncFeatureBranchSteps(branchName string) (result StepList) {
	if git.HasTrackingBranch(branchName) {
		result.Append(&MergeBranchStep{BranchName: git.GetTrackingBranchName(branchName)})
	}
	result.Append(&MergeBranchStep{BranchName: git.Config().GetParentBranch(branchName)})
	return
}

func getSyncNonFeatureBranchSteps(branchName string) (result StepList) {
	if git.HasTrackingBranch(branchName) {
		if git.Config().GetPullBranchStrategy() == "rebase" {
			result.Append(&RebaseBranchStep{BranchName: git.GetTrackingBranchName(branchName)})
		} else {
			result.Append(&MergeBranchStep{BranchName: git.GetTrackingBranchName(branchName)})
		}
	}

	mainBranchName := git.Config().GetMainBranch()
	if mainBranchName == branchName && git.HasRemote("upstream") && git.Config().ShouldSyncUpstream() {
		result.Append(&FetchUpstreamStep{BranchName: mainBranchName})
		result.Append(&RebaseBranchStep{BranchName: fmt.Sprintf("upstream/%s", mainBranchName)})
	}
	return
}
