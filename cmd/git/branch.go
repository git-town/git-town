package git

import (
  "fmt"
  "io/ioutil"
  "log"
  "strings"

  "github.com/Originate/gt/cmd/util"
)


func EnsureDoesNotHaveBranch(branchName string) {
  if HasBranch(branchName) {
    util.ExitWithErrorMessage(fmt.Sprintf("A branch named '%s' already exists", branchName))
  }
}


func GetCurrentBranchName() string {
  if IsRebaseInProgress() {
    filename := fmt.Sprintf("%s/.git/rebase-apply/head-name", GetRootDirectory())
    rawContent, err := ioutil.ReadFile(filename)
    if err != nil {
      log.Fatal(err)
    }
    content := strings.TrimSpace(string(rawContent))
    return strings.Replace(content, "refs/heads/", "", -1)
  } else {
    return util.GetCommandOutput([]string{"git", "rev-parse", "--abbrev-ref", "HEAD"})
  }
}


func GetTrackingBranchName(branchName string) string {
  return "origin/" + branchName
}


func HasBranch(branchName string) bool {
  cmd := []string{"git", "branch", "-a"}
  lines := strings.Split(util.GetCommandOutput(cmd), "\n")
  for _, line := range(lines) {
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
  output := util.GetCommandOutput([]string{"git", "branch", "-r"})
  lines := strings.Split(output, "\n")
  for i := 0; i < len(lines); i++ {
    line := lines[i]
    line = strings.TrimSpace(line)
    if line == trackingBranchName {
      return true
    }
  }
  return false
}


func ShouldBranchBePushed(branchName string) bool {
  trackingBranchName := GetTrackingBranchName(branchName)
  output := util.GetCommandOutput([]string{"git", "rev-list", "--left-right", branchName + "..." + trackingBranchName})
  return output != ""
}
