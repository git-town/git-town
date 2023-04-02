package execute

import (
	"github.com/git-town/git-town/v7/src/cache"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/subshell"
	"github.com/git-town/git-town/v7/src/validate"
)

func LoadProdRunner(args LoadArgs) (prodRunner git.ProdRunner, exit bool, err error) { //nolint:nonamedreturns // so many return values require names
	var stats *Statistics
	if args.Debug {
		stats = &Statistics{}
	}
	backendRunner := NewBackendRunner(nil, args.Debug, stats)
	config := git.NewRepoConfig(backendRunner)
	frontendRunner := NewFrontendRunner(args.OmitBranchNames, args.DryRun, config.CurrentBranchCache, stats)
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
	ec := runstate.ErrorChecker{}
	if args.ValidateGitversion {
		ec.Check(validate.HasGitVersion(&prodRunner.Backend))
	}
	if args.ValidateIsConfigured {
		ec.Check(validate.IsConfigured(&prodRunner.Backend))
	}
	if args.ValidateIsOnline {
		ec.Check(validate.IsOnline(&prodRunner.Config))
	}
	if args.HandleUnfinishedState {
		exit = ec.Bool(validate.HandleUnfinishedState(&prodRunner, nil))
	}
	return prodRunner, exit, ec.Err
}

type LoadArgs struct {
	Debug                 bool
	DryRun                bool
	HandleUnfinishedState bool
	OmitBranchNames       bool `exhaustruct:"optional"`
	ValidateGitversion    bool `exhaustruct:"optional"`
	ValidateIsRepository  bool `exhaustruct:"optional"`
	ValidateIsConfigured  bool `exhaustruct:"optional"`
	ValidateIsOnline      bool `exhaustruct:"optional"`
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
