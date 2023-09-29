package execute

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/src/undo"
	"github.com/git-town/git-town/v9/src/validate"
)

// LoadBranches loads the typically used information about Git branches using a single Git command.
func LoadBranches(args LoadBranchesArgs) (domain.Branches, domain.BranchesSnapshot, undo.StashSnapshot, bool, error) {
	var branchesSnapshot domain.BranchesSnapshot
	var err error
	stashSize, err := args.Repo.Runner.Backend.StashSize()
	if err != nil {
		return domain.EmptyBranches(), domain.EmptyBranchesSnapshot(), undo.EmptyStashSnapshot(), false, err
	}
	stashSnapshot := undo.StashSnapshot{
		Amount: stashSize,
	}
	if args.HandleUnfinishedState {
		branchesSnapshot, err = args.Repo.Runner.Backend.BranchesSnapshot()
		if err != nil {
			return domain.EmptyBranches(), domain.EmptyBranchesSnapshot(), undo.EmptyStashSnapshot(), false, err
		}
		exit, err := validate.HandleUnfinishedState(&args.Repo.Runner, nil, args.Repo.RootDir, args.Lineage, branchesSnapshot, args.Repo.ConfigSnapshot, stashSnapshot, args.PushHook)
		if err != nil || exit {
			return domain.EmptyBranches(), domain.EmptyBranchesSnapshot(), undo.EmptyStashSnapshot(), exit, err
		}
	}
	if args.ValidateNoOpenChanges {
		hasOpenChanges, err := args.Repo.Runner.Backend.HasOpenChanges()
		if err != nil {
			return domain.EmptyBranches(), domain.EmptyBranchesSnapshot(), undo.EmptyStashSnapshot(), false, err
		}
		err = validate.NoOpenChanges(hasOpenChanges)
		if err != nil {
			return domain.EmptyBranches(), domain.EmptyBranchesSnapshot(), undo.EmptyStashSnapshot(), false, err
		}
	}
	if args.Fetch {
		var remotes domain.Remotes
		remotes, err := args.Repo.Runner.Backend.Remotes()
		if err != nil {
			return domain.EmptyBranches(), domain.EmptyBranchesSnapshot(), undo.EmptyStashSnapshot(), false, err
		}
		if remotes.HasOrigin() && !args.Repo.IsOffline {
			err = args.Repo.Runner.Frontend.Fetch()
			if err != nil {
				return domain.EmptyBranches(), domain.EmptyBranchesSnapshot(), undo.EmptyStashSnapshot(), false, err
			}
		}
		branchesSnapshot, err = args.Repo.Runner.Backend.BranchesSnapshot()
		if err != nil {
			return domain.EmptyBranches(), domain.EmptyBranchesSnapshot(), undo.EmptyStashSnapshot(), false, err
		}
	}
	if branchesSnapshot.IsEmpty() {
		branchesSnapshot, err = args.Repo.Runner.Backend.BranchesSnapshot()
		if err != nil {
			return domain.EmptyBranches(), domain.EmptyBranchesSnapshot(), undo.EmptyStashSnapshot(), false, err
		}
	}
	branchTypes := args.Repo.Runner.Config.BranchTypes()
	branches := domain.Branches{
		All:     branchesSnapshot.Branches,
		Types:   branchTypes,
		Initial: branchesSnapshot.Active,
	}
	if args.ValidateIsConfigured {
		branches.Types, err = validate.IsConfigured(&args.Repo.Runner.Backend, branches)
	}
	return branches, branchesSnapshot, stashSnapshot, false, err
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
