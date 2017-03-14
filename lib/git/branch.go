package git

import (
  "fmt"
  "io/ioutil"
  "log"
  "strings"

  "github.com/Originate/git-town/lib/util"
)


func EnsureDoesNotHaveBranch(branchName string) {
  if HasBranch(branchName) {
    util.ExitWithErrorMessage(fmt.Sprintf("A branch named '%s' already exists", branchName))
  }
}


func GetCurrentBranchName() string {
  if IsRebaseInProgress() {
    return getCurrentBranchNameDuringRebase()
  } else {
    return util.GetCommandOutput("git", "rev-parse", "--abbrev-ref", "HEAD")
  }
}


func GetTrackingBranchName(branchName string) string {
  return "origin/" + branchName
}


func HasBranch(branchName string) bool {
  output := util.GetCommandOutput("git", "branch", "-a")
  for _, line := range(strings.Split(output, "\n")) {
    line = strings.Trim(line, "* ")
    line = strings.TrimSpace(line)
    line = strings.Replace(line, "remotes/origin/", "", 1)
    if line == branchName {
      return true
    }
  }
  return false
}


func HasTrackingBranch(branchName string) bool {
  trackingBranchName := GetTrackingBranchName(branchName)
  output := util.GetCommandOutput("git", "branch", "-r")
  for _, line := range(strings.Split(output, "\n")) {
    if strings.TrimSpace(line) == trackingBranchName {
      return true
    }
  }
  return false
}


func ShouldBranchBePushed(branchName string) bool {
  trackingBranchName := GetTrackingBranchName(branchName)
  output := util.GetCommandOutput("git", "rev-list", "--left-right", branchName + "..." + trackingBranchName)
  return output != ""
}


// Helpers


func getCurrentBranchNameDuringRebase() string {
  filename := fmt.Sprintf("%s/.git/rebase-apply/head-name", GetRootDirectory())
  rawContent, err := ioutil.ReadFile(filename)
  if err != nil {
    log.Fatal(err)
  }
  content := strings.TrimSpace(string(rawContent))
  return strings.Replace(content, "refs/heads/", "", -1)
}
