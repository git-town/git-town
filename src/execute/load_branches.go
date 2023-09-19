package execute

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/undo"
	"github.com/git-town/git-town/v9/src/validate"
)

// LoadBranches loads the typically used information about Git branches using a single Git command.
func LoadBranches(args LoadBranchesArgs) (domain.Branches, undo.BranchesSnapshot, bool, error) {
	var allBranches domain.BranchInfos
	var initialBranch domain.LocalBranchName
	var err error
	if args.HandleUnfinishedState {
		allBranches, initialBranch, err = args.Repo.Runner.Backend.BranchInfos()
		if err != nil {
			return domain.EmptyBranches(), undo.EmptyBranchesSnapshot(), false, err
		}
		initialBranchesSnapshot := undo.BranchesSnapshot{
			Branches: allBranches,
			Active:   initialBranch,
		}
		exit, err := validate.HandleUnfinishedState(&args.Repo.Runner, nil, args.Repo.RootDir, args.Lineage, initialBranchesSnapshot, args.Repo.ConfigSnapshot)
		if err != nil || exit {
			return domain.EmptyBranches(), undo.EmptyBranchesSnapshot(), exit, err
		}
	}
	if args.ValidateNoOpenChanges {
		hasOpenChanges, err := args.Repo.Runner.Backend.HasOpenChanges()
		if err != nil {
			return domain.EmptyBranches(), undo.EmptyBranchesSnapshot(), false, err
		}
		err = validate.NoOpenChanges(hasOpenChanges)
		if err != nil {
			return domain.EmptyBranches(), undo.EmptyBranchesSnapshot(), false, err
		}
	}
	if args.Fetch {
		var remotes domain.Remotes
		remotes, err := args.Repo.Runner.Backend.Remotes()
		if err != nil {
			return domain.EmptyBranches(), undo.EmptyBranchesSnapshot(), false, err
		}
		if remotes.HasOrigin() && !args.Repo.IsOffline {
			err = args.Repo.Runner.Frontend.Fetch()
			if err != nil {
				return domain.EmptyBranches(), undo.EmptyBranchesSnapshot(), false, err
			}
		}
		allBranches, initialBranch, err = args.Repo.Runner.Backend.BranchInfos()
		if err != nil {
			return domain.EmptyBranches(), undo.EmptyBranchesSnapshot(), false, err
		}
	}
	if initialBranch.IsEmpty() {
		allBranches, initialBranch, err = args.Repo.Runner.Backend.BranchInfos()
		if err != nil {
			return domain.EmptyBranches(), undo.EmptyBranchesSnapshot(), false, err
		}
	}
	branchTypes := args.Repo.Runner.Config.BranchTypes()
	branches := domain.Branches{
		All:     allBranches,
		Types:   branchTypes,
		Initial: initialBranch,
	}
	if args.ValidateIsConfigured {
		branches.Types, err = validate.IsConfigured(&args.Repo.Runner.Backend, branches)
	}
	snapshot := undo.BranchesSnapshot{
		Branches: branches.All.Clone(),
		Active:   initialBranch,
	}
	return branches, snapshot, false, err
}

type LoadBranchesArgs struct {
	Repo                  *OpenRepoResult
	Fetch                 bool
	HandleUnfinishedState bool
	Lineage               config.Lineage
	ValidateIsConfigured  bool
	ValidateNoOpenChanges bool
}
