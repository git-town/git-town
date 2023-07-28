package execute

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/validate"
)

// LoadBranches loads the typically used information about Git branches using a single Git command.
func LoadBranches(pr *git.ProdRunner, args LoadBranchesArgs) (Branches, error) {
	allBranches, initialBranch, err := pr.Backend.BranchesSyncStatus()
	if err != nil {
		return EmptyBranches(), err
	}
	branchDurations := pr.Config.BranchDurations()
	if args.ValidateIsConfigured {
		branchDurations, err = validate.IsConfigured(&pr.Backend, allBranches, branchDurations)
	}
	return Branches{
		All:       allBranches,
		Durations: branchDurations,
		Initial:   initialBranch,
	}, err
}

type LoadBranchesArgs struct {
	ValidateIsConfigured bool
}

type Branches struct {
	All       git.BranchesSyncStatus
	Durations config.BranchDurations
	Initial   string
}

// EmptyBranches provides the zero value for Branches.
func EmptyBranches() Branches {
	return Branches{
		All:       git.BranchesSyncStatus{},
		Durations: config.EmptyBranchDurations(),
		Initial:   "",
	}
}
