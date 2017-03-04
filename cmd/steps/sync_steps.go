package steps

import (
  "fmt"

  "github.com/Originate/gt/cmd/config"
  "github.com/Originate/gt/cmd/git"
)


func GetSyncBranchSteps(branchName string) []Step {
  isFeature := config.IsFeatureBranch(branchName)
  hasRemoteOrigin := config.HasRemoteOrigin()

  var steps []Step

  if hasRemoteOrigin || isFeature {
    steps = append(steps, CheckoutBranchStep{BranchName: branchName})
    if isFeature {
      steps = append(steps, MergeTrackingBranchStep{}, MergeBranchStep{BranchName: config.GetParentBranch(branchName)})
    } else {
      if config.GetPullBranchStrategy() == "rebase" {
        steps = append(steps, RebaseTrackingBranchStep{})
      } else {
        steps = append(steps, MergeTrackingBranchStep{})
      }

      mainBranchName := config.GetMainBranch()
      if mainBranchName == branchName && config.HasRemoteUpstream() {
        steps = append(steps, FetchUpstreamStep{}, RebaseBranchStep{BranchName: fmt.Sprintf("upstream/%s", mainBranchName)})
      }
    }

    if hasRemoteOrigin {
      if git.HasTrackingBranch(branchName) {
        steps = append(steps, PushBranchStep{BranchName: branchName})
      } else {
        steps = append(steps, new(CreateTrackingBranchStep))
      }
    }
  }

  return steps
}
