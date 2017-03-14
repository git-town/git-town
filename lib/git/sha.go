package git

import (
  "github.com/Originate/git-town/lib/util"
)


func GetBranchSha(branchName string) string {
  return util.GetCommandOutput([]string{"git", "rev-parse", branchName})
}


func GetCurrentSha() string {
  return GetBranchSha("HEAD")
}
