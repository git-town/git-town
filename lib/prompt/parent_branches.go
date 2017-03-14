package prompt

import (
  "fmt"

  "github.com/Originate/gt/lib/config"
)


func EnsureKnowsParentBranches(branchNames []string) {
  headerShown := false
  for _, branchName := range(branchNames) {
    if config.IsMainBranch(branchName) || config.IsPerennialBranch(branchName) || config.HasStoredAncestorBranches(branchName) {
      continue
    }
    if !headerShown {
      printParentBranchHeader()
      headerShown = true
    }
    askForBranchAncestry(branchName)
    ancestors := config.CompileAncestorBranches(branchName)
    config.SetAncestorBranches(branchName, ancestors)
  }
  if headerShown {
    fmt.Println()
  }
}


// Helpers

func askForBranchAncestry(branchName string) {
  current := branchName
  for {
    parent := config.GetParentBranch(current)
    if parent == "" {
      parent = askForParentBranch(current)
      config.SetParentBranch(current, parent)
    }
    if parent == config.GetMainBranch() || config.IsPerennialBranch(parent) {
      break
    }
    current = parent
  }
}

func askForParentBranch(branchName string) string {
  panic("unimplemented")
}

func printParentBranchHeader() {
  // TODO
}
