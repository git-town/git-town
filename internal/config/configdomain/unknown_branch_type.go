package configdomain

import . "github.com/git-town/git-town/v22/pkg/prelude"

// UnknownBranchType is the type that branches downloaded from the dev remote have by default.
type UnknownBranchType BranchType

func (self UnknownBranchType) BranchType() BranchType {
	return BranchType(self)
}

func (self UnknownBranchType) String() string {
	return self.BranchType().String()
}

func UnknownBranchTypeOpt(branchTypeOpt Option[BranchType]) Option[UnknownBranchType] {
	if branchType, has := branchTypeOpt.Get(); has {
		return Some(UnknownBranchType(branchType))
	}
	return None[UnknownBranchType]()
}
