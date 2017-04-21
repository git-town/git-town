package steps

import (
	"fmt"

	"github.com/Originate/git-town/lib/git"
)

// GetSyncBranchSteps returns the steps to sync the branch with the given name.
func GetSyncBranchSteps(branchName string) (result StepList) {
	isFeature := git.IsFeatureBranch(branchName)
	hasRemoteOrigin := git.HasRemote("origin")

	if !hasRemoteOrigin && !isFeature {
		return
	}

	result.Append(CheckoutBranchStep{BranchName: branchName})
	if isFeature {
		result.Append(MergeTrackingBranchStep{})
		result.Append(MergeBranchStep{BranchName: git.GetParentBranch(branchName)})
	} else {
		if git.GetPullBranchStrategy() == "rebase" {
			result.Append(RebaseTrackingBranchStep{})
		} else {
			result.Append(MergeTrackingBranchStep{})
		}

		mainBranchName := git.GetMainBranch()
		if mainBranchName == branchName && git.HasRemote("upstream") {
			result.Append(FetchUpstreamStep{})
			result.Append(RebaseBranchStep{BranchName: fmt.Sprintf("upstream/%s", mainBranchName)})
		}
	}

	if hasRemoteOrigin {
		if git.HasTrackingBranch(branchName) {
			result.Append(PushBranchStep{BranchName: branchName})
		} else {
			result.Append(CreateTrackingBranchStep{BranchName: branchName})
		}
	}

	return
}
