package steps

import (
  "fmt"

  "github.com/Originate/git-town/lib/config"
  "github.com/Originate/git-town/lib/git"
)


func GetSyncBranchSteps(branchName string) (result StepList) {
  isFeature := config.IsFeatureBranch(branchName)
  hasRemoteOrigin := config.HasRemote("origin")

  if !hasRemoteOrigin && !isFeature {
    return
  }

  result.Append(CheckoutBranchStep{BranchName: branchName})
  if isFeature {
    result.Append(MergeTrackingBranchStep{})
    result.Append(MergeBranchStep{BranchName: config.GetParentBranch(branchName)})
  } else {
    if config.GetPullBranchStrategy() == "rebase" {
      result.Append(RebaseTrackingBranchStep{})
    } else {
      result.Append(MergeTrackingBranchStep{})
    }

    mainBranchName := config.GetMainBranch()
    if mainBranchName == branchName && config.HasRemote("upstream") {
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
