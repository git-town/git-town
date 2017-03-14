package git

import (
  "github.com/Originate/gt/lib/util"
)


func GetBranchSha(branchName string) string {
  return util.GetCommandOutput("git", "rev-parse", branchName)
}


func GetCurrentSha() string {
  return GetBranchSha("HEAD")
}
