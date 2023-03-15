package subshell

type CurrentBranchTracker struct {
	// the current branch
	Value string
}

func (cbt *CurrentBranchTracker) Track(executable string, args ...string) {
	if executable == "git" && args[0] == "checkout" && len(args) == 2 {
		cbt.Value = args[1]
	}
}
