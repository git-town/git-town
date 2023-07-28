package hosting_test

// emptyShaForBranch is a dummy implementation for hosting.ShaForBranchfunc to be used in tests.
func emptyShaForBranch(string) (string, error) {
	return "", nil
}
