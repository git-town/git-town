package configdomain

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

// BranchTypes answers questions about whether branches are long-lived or not.
type BranchTypes struct {
	MainBranch        gitdomain.LocalBranchName
	PerennialBranches gitdomain.LocalBranchNames
}

func EmptyBranchTypes() BranchTypes {
	return BranchTypes{
		MainBranch:        gitdomain.EmptyLocalBranchName(),
		PerennialBranches: gitdomain.LocalBranchNames{},
	}
}
