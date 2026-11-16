package configdomain

import . "github.com/git-town/git-town/v24/pkg/prelude"

// NewBranchType is the type that branches created with hack, append, and prepend have.
type NewBranchType BranchType

func NewBranchTypeOpt(value Option[BranchType]) Option[NewBranchType] {
	if branchType, has := value.Get(); has {
		return Some(NewBranchType(branchType))
	}
	return None[NewBranchType]()
}

func (self NewBranchType) BranchType() BranchType {
	return BranchType(self)
}

func (self NewBranchType) String() string {
	return self.BranchType().String()
}
