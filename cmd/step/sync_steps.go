package step

import (
  "github.com/Originate/gt/cmd/config"
  "github.com/Originate/gt/cmd/git"
)

func GetSyncBranchSteps(branchName string) []Step {
  isFeature := config.IsFeatureBranch(branchName)
  hasRemote := config.HasRemote()

  var steps []Step

  if hasRemote || isFeature {
    steps = append(steps, CheckoutBranchStep{BranchName: branchName})
    if isFeature {
      steps = append(steps, new(MergeTrackingBranchStep), MergeBranchStep{BranchName: config.GetParentBranch(branchName)})
    } else {
      if config.GetPullBranchStrategy() == "rebase" {
        steps = append(steps, new(RebaseTrackingBranchStep))
      } else {
        steps = append(steps, new(MergeTrackingBranchStep))
      }

      // if config.getMainBranchName() == branchName && git.hasRemoteUpstream() {
      //   append(steps, new(FetchUpstreamStep))
      //   append(steps, new(RebaseBranchStep(branchName: fmt.Sprintf("upstream/%s", config.mainBranchName))))
      // }
    }

    if hasRemote {
      if git.HasTrackingBranch(branchName) {
        steps = append(steps, PushBranchStep{BranchName: branchName})
      } else {
        steps = append(steps, new(CreateTrackingBranchStep))
      }
    }
  }

  return steps
}
