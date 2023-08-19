package execute

import (
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/validate"
)

// LoadBranches loads the typically used information about Git branches using a single Git command.
func LoadBranches(pr *git.ProdRunner, args LoadBranchesArgs) (domain.Branches, error) {
	allBranches, initialBranch, err := pr.Backend.BranchInfos()
	if err != nil {
		return domain.EmptyBranches(), err
	}
	branchDurations := pr.Config.BranchTypes()
	if args.ValidateIsConfigured {
		branchDurations, err = validate.IsConfigured(&pr.Backend, allBranches, branchDurations)
	}
	return domain.Branches{
		All:        allBranches,
		Perennials: branchDurations,
		Initial:    initialBranch,
	}, err
}

type LoadBranchesArgs struct {
	ValidateIsConfigured bool
}
