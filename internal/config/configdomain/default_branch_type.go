package configdomain

import . "github.com/git-town/git-town/v16/pkg/prelude"

// the configured default branch type
type DefaultBranchType struct {
	BranchType
}

func ParseDefaultBranchType(text string) (Option[DefaultBranchType], error) {
	branchTypeOpt, err := ParseBranchType(text)
	if branchType, hasBranchType := branchTypeOpt.Get(); hasBranchType {
		return Some(DefaultBranchType{BranchType: branchType}), err
	}
	return None[DefaultBranchType](), err
}
