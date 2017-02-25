package step

import (
  "github.com/Originate/gt/cmd/git"
  "github.com/Originate/gt/cmd/util"
)

func GetSyncBranchSteps(branchName string, config util.Config) []Step {
  isFeature := git.IsFeatureBranch(branchName)

  var steps []Step

  if config.HasRemote || isFeature {
    steps = append(steps, CheckoutBranchStep{branchName: branchName})
    if isFeature {
      steps = append(steps, new(MergeTrackingBranchStep), MergeBranchStep{branchName: git.GetParentBranch(branchName)})
    } else {
      if config.PullBranchStrategy == "rebase" {
        steps = append(steps, new(RebaseTrackingBranchStep))
      } else {
        steps = append(steps, new(MergeTrackingBranchStep))
      }

      // if config.mainBranchName == branchName && git.hasRemoteUpstream() {
      //   append(steps, new(FetchUpstreamStep))
      //   append(steps, new(RebaseBranchStep(branchName: fmt.Sprintf("upstream/%s", config.mainBranchName))))
      // }
    }

    if config.HasRemote {
      if git.HasTrackingBranch(branchName) {
        steps = append(steps, PushBranchStep{branchName: branchName})
      } else {
        steps = append(steps, new(CreateTrackingBranchStep))
      }
    }
  }

  return steps
}
