package execute

import (
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/validate"
)

func LoadBranches(pr *git.ProdRunner, args LoadBranchesArgs) (allBranches git.BranchesSyncStatus, currentBranch string, err error) {
	allBranches, currentBranch, err = pr.Backend.BranchesSyncStatus()
	if err != nil {
		return
	}
	if args.ValidateIsConfigured {
		err = validate.IsConfigured(&pr.Backend, allBranches)
		if err != nil {
			return
		}
	}
	return
}

type LoadBranchesArgs struct {
	ValidateIsConfigured bool
}
