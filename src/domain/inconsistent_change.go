package domain

// InconsistentChange describes a change where both local and remote branch exist before and after,
// but it's not an OmniChange, i.e. the SHA are different.
type InconsistentChange struct {
	Before BranchInfo
	After  BranchInfo
}
