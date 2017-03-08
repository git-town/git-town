package prompt


func EnsureKnowsParentBranches(branchNames []string) {
  for _, branchName := range(branchNames) {
    EnsureKnowsParentBranch(branchName)
  }
}

func EnsureKnowsParentBranch(branchName string) {

}
