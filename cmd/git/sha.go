package git

import (
  "github.com/Originate/gt/cmd/util"
)


func GetCurrentSha() string {
  return GetBranchSha("HEAD")
}


func GetBranchSha(branchName string) string {
  return util.GetCommandOutput([]string{"git", "rev-parse", branchName})
}
