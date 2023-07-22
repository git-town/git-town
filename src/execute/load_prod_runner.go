package execute

import (
	"github.com/git-town/git-town/v9/src/cache"
	"github.com/git-town/git-town/v9/src/failure"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/statistics"
	"github.com/git-town/git-town/v9/src/subshell"
	"github.com/git-town/git-town/v9/src/validate"
)

func LoadProdRunner(args LoadArgs) (prodRunner git.ProdRunner, exit bool, err error) { //nolint:nonamedreturns // so many return values require names
	var stats Statistics
	if args.Debug {
		stats = &statistics.CommandsRun{CommandsCount: 0}
	} else {
		stats = &statistics.None{}
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
	if args.ValidateIsRepository {
		err := validate.IsRepository(&prodRunner)
		if err != nil {
			return prodRunner, false, err
		}
	}
	if !args.OmitBranchNames || args.DryRun {
		currentBranch, err := prodRunner.Backend.CurrentBranch()
		if err != nil {
			return prodRunner, false, err
		}
		prodRunner.Config.CurrentBranchCache.Set(currentBranch)
	}
	if args.DryRun {
		prodRunner.Config.DryRun = true
	}
	fc := failure.Collector{}
	if args.ValidateGitversion {
		fc.Check(validate.HasGitVersion(&prodRunner.Backend))
	}
	if args.ValidateIsConfigured {
		fc.Check(validate.IsConfigured(&prodRunner.Backend))
	}
	if args.ValidateIsOnline {
		fc.Check(validate.IsOnline(&prodRunner.Config))
	}
	if args.HandleUnfinishedState {
		exit = fc.Bool(validate.HandleUnfinishedState(&prodRunner, nil))
	}
	return prodRunner, exit, fc.Err
}

type LoadArgs struct {
	Debug                 bool
	DryRun                bool
	HandleUnfinishedState bool
	OmitBranchNames       bool
	ValidateGitversion    bool
	ValidateIsRepository  bool
	ValidateIsConfigured  bool
	ValidateIsOnline      bool
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
