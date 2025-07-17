package configdomain

import . "github.com/git-town/git-town/v21/pkg/prelude"

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
