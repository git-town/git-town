package execute

import (
	"github.com/git-town/git-town/v9/src/cache"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/subshell"
	"github.com/git-town/git-town/v9/src/validate"
)

func LoadProdRunner(args LoadArgs) (prodRunner git.ProdRunner, err error) { //nolint:nonamedreturns // so many return values require names
	var stats Statistics
	if args.Debug {
		stats = &CommandsStatistics{CommandsCount: 0}
	} else {
		stats = &NoStatistics{}
	}
	backendRunner := subshell.BackendRunner{Dir: nil, Verbose: args.Debug, Stats: stats}
	config := git.NewRepoConfig(backendRunner)
	prodRunner = git.ProdRunner{
		Config: config,
		Backend: git.BackendCommands{
			BackendRunner: backendRunner,
			Config:        &config,
		},
		Frontend: git.FrontendCommands{
			FrontendRunner: NewFrontendRunner(args.OmitBranchNames, args.DryRun, config.CurrentBranchCache, stats),
			Config:         &config,
		},
		Stats: stats,
	}
	err = validate.HasGitVersion(&prodRunner.Backend)
	if err != nil {
		return prodRunner, err
	}
	if args.ValidateIsConfigured {
		err = validate.IsConfigured(&prodRunner.Backend)
		if err != nil {
			return prodRunner, err
		}
	}
	if !args.OmitBranchNames || args.DryRun {
		currentBranch, err := prodRunner.Backend.CurrentBranch()
		if err != nil {
			return prodRunner, err
		}
		prodRunner.Config.CurrentBranchCache.Set(currentBranch)
	}
	if args.DryRun {
		prodRunner.Config.DryRun = true
	}
	return prodRunner, nil
}

type LoadArgs struct {
	Debug                bool
	DryRun               bool
	OmitBranchNames      bool
	ValidateIsConfigured bool
}

// NewFrontendRunner provides a FrontendRunner instance that behaves according to the given configuration.
func NewFrontendRunner(omitBranchNames, dryRun bool, currentBranchCache *cache.String, stats Statistics) git.FrontendRunner {
	if dryRun {
		return &subshell.FrontendDryRunner{
			CurrentBranch:   currentBranchCache,
			OmitBranchNames: omitBranchNames,
			Stats:           stats,
		}
	}
	return &subshell.FrontendRunner{
		CurrentBranch:   currentBranchCache,
		OmitBranchNames: omitBranchNames,
		Stats:           stats,
	}
}
