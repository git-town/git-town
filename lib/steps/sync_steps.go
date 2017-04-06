package steps

import (
	"fmt"

	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/gitconfig"
)

func GetSyncBranchSteps(branchName string) (result StepList) {
	isFeature := gitconfig.IsFeatureBranch(branchName)
	hasRemoteOrigin := gitconfig.HasRemote("origin")

	if !hasRemoteOrigin && !isFeature {
		return
	}

	result.Append(CheckoutBranchStep{BranchName: branchName})
	if isFeature {
		result.Append(MergeTrackingBranchStep{})
		result.Append(MergeBranchStep{BranchName: gitconfig.GetParentBranch(branchName)})
	} else {
		if gitconfig.GetPullBranchStrategy() == "rebase" {
			result.Append(RebaseTrackingBranchStep{})
		} else {
			result.Append(MergeTrackingBranchStep{})
		}

		mainBranchName := gitconfig.GetMainBranch()
		if mainBranchName == branchName && gitconfig.HasRemote("upstream") {
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
