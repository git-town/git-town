package execute

import (
	"github.com/git-town/git-town/v7/src/cache"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/subshell"
	"github.com/git-town/git-town/v7/src/validate"
)

func LoadProdRunner(args loadArgs) (prodRunner git.ProdRunner, exit bool, err error) { //nolint:nonamedreturns // so many return values require names
	var stats *Statistics
	if args.debug {
		stats = &Statistics{}
	}
	backendRunner := NewBackendRunner(nil, args.debug, stats)
	config := git.NewRepoConfig(backendRunner)
	frontendRunner := NewFrontendRunner(args.omitBranchNames, args.dryRun, config.CurrentBranchCache, stats)
	backendCommands := git.BackendCommands{
		BackendRunner: backendRunner,
		Config:        &config,
	}
	prodRunner = git.ProdRunner{
		Config:  config,
		Backend: backendCommands,
		Frontend: git.FrontendCommands{
			Frontend: frontendRunner,
			Config:   &config,
			Backend:  &backendCommands,
		},
		Stats: stats,
	}
	if args.validateIsRepository {
		err := validate.IsRepository(&prodRunner)
		if err != nil {
			return prodRunner, false, err
		}
	}
	if !args.omitBranchNames || args.dryRun {
		currentBranch, err := prodRunner.Backend.CurrentBranch()
		if err != nil {
			return prodRunner, false, err
		}
		prodRunner.Config.CurrentBranchCache.Set(currentBranch)
	}
	if args.dryRun {
		prodRunner.Config.DryRun = true
	}
	ec := runstate.ErrorChecker{}
	if args.validateGitversion {
		ec.Check(validate.HasGitVersion(&prodRunner.Backend))
	}
	if args.validateIsConfigured {
		ec.Check(validate.IsConfigured(&prodRunner.Backend))
	}
	if args.validateIsOnline {
		ec.Check(validate.IsOnline(&prodRunner.Config))
	}
	if args.handleUnfinishedState {
		exit = ec.Bool(validate.HandleUnfinishedState(&prodRunner, nil))
	}
	return prodRunner, exit, ec.Err
}

type loadArgs struct {
	debug                 bool
	dryRun                bool
	handleUnfinishedState bool
	omitBranchNames       bool `exhaustruct:"optional"`
	validateGitversion    bool `exhaustruct:"optional"`
	validateIsRepository  bool `exhaustruct:"optional"`
	validateIsConfigured  bool `exhaustruct:"optional"`
	validateIsOnline      bool `exhaustruct:"optional"`
}

func NewBackendRunner(dir *string, debug bool, statistics *Statistics) git.BackendRunner {
	backendRunner := subshell.BackendRunner{Dir: dir, Statistics: statistics}
	if debug {
		return subshell.BackendLoggingRunner{Runner: backendRunner, Statistics: statistics}
	}
	return backendRunner
}

// NewFrontendRunner provides a FrontendRunner instance that behaves according to the given configuration.
func NewFrontendRunner(omitBranchNames, dryRun bool, currentBranchCache *cache.String, stats *Statistics) git.FrontendRunner {
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
