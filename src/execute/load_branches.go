package execute

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/undo"
	"github.com/git-town/git-town/v9/src/validate"
)

// LoadBranches loads the typically used information about Git branches using a single Git command.
func LoadBranches(args LoadBranchesArgs) (domain.Branches, undo.BranchesSnapshot, undo.StashSnapshot, bool, error) {
	var allBranches domain.BranchInfos
	var initialBranch domain.LocalBranchName
	var err error
	stashSize, err := args.Repo.Runner.Backend.StashSize()
	if err != nil {
		return domain.EmptyBranches(), undo.EmptyBranchesSnapshot(), undo.EmptyStashSnapshot(), false, err
	}
	initialStashSnapshot := undo.StashSnapshot{
		Amount: stashSize,
	}
	if args.HandleUnfinishedState {
		allBranches, initialBranch, err = args.Repo.Runner.Backend.BranchInfos()
		if err != nil {
			return domain.EmptyBranches(), undo.EmptyBranchesSnapshot(), undo.EmptyStashSnapshot(), false, err
		}
		initialBranchesSnapshot := undo.BranchesSnapshot{
			Branches: allBranches,
			Active:   initialBranch,
		}
		exit, err := validate.HandleUnfinishedState(&args.Repo.Runner, nil, args.Repo.RootDir, args.Lineage, initialBranchesSnapshot, args.Repo.ConfigSnapshot, initialStashSnapshot, args.PushHook)
		if err != nil || exit {
			return domain.EmptyBranches(), undo.EmptyBranchesSnapshot(), undo.EmptyStashSnapshot(), exit, err
		}
	}
	if args.ValidateNoOpenChanges {
		hasOpenChanges, err := args.Repo.Runner.Backend.HasOpenChanges()
		if err != nil {
			return domain.EmptyBranches(), undo.EmptyBranchesSnapshot(), undo.EmptyStashSnapshot(), false, err
		}
		err = validate.NoOpenChanges(hasOpenChanges)
		if err != nil {
			return domain.EmptyBranches(), undo.EmptyBranchesSnapshot(), undo.EmptyStashSnapshot(), false, err
		}
	}
	if args.Fetch {
		var remotes domain.Remotes
		remotes, err := args.Repo.Runner.Backend.Remotes()
		if err != nil {
			return domain.EmptyBranches(), undo.EmptyBranchesSnapshot(), undo.EmptyStashSnapshot(), false, err
		}
		if remotes.HasOrigin() && !args.Repo.IsOffline {
			err = args.Repo.Runner.Frontend.Fetch()
			if err != nil {
				return domain.EmptyBranches(), undo.EmptyBranchesSnapshot(), undo.EmptyStashSnapshot(), false, err
			}
		}
		allBranches, initialBranch, err = args.Repo.Runner.Backend.BranchInfos()
		if err != nil {
			return domain.EmptyBranches(), undo.EmptyBranchesSnapshot(), undo.EmptyStashSnapshot(), false, err
		}
	}
	if initialBranch.IsEmpty() {
		allBranches, initialBranch, err = args.Repo.Runner.Backend.BranchInfos()
		if err != nil {
			return domain.EmptyBranches(), undo.EmptyBranchesSnapshot(), undo.EmptyStashSnapshot(), false, err
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
	return branches, snapshot, initialStashSnapshot, false, err
}

type LoadBranchesArgs struct {
	Repo                  *OpenRepoResult
	Fetch                 bool
	HandleUnfinishedState bool
	Lineage               config.Lineage
	PushHook              bool
	ValidateIsConfigured  bool
	ValidateNoOpenChanges bool
}
