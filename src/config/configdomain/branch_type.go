package configdomain

type BranchType int

const (
	BranchTypeMainBranch BranchType = iota
	BranchTypePerennialBranch
	BranchTypeFeatureBranch
	BranchTypeObservedBranch
	BranchTypeParkedBranch
)

// ShouldPush indicates whether a branch with this type should push its local commit to origin.
func (self BranchType) ShouldPush() bool {
	switch self {
	case BranchTypeMainBranch, BranchTypeFeatureBranch, BranchTypePerennialBranch, BranchTypeParkedBranch:
		return true
	case BranchTypeObservedBranch:
		return false
	}
	panic("unhandled branch type")
}
