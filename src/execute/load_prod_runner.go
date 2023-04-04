package execute

import (
	"github.com/git-town/git-town/v7/src/cache"
	"github.com/git-town/git-town/v7/src/failure"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/subshell"
	"github.com/git-town/git-town/v7/src/validate"
)

func LoadProdRunner(args LoadArgs) (prodRunner git.ProdRunner, exit bool, err error) { //nolint:nonamedreturns // so many return values require names
	backendRunner := NewBackendRunner(nil, args.Debug)
	config := git.NewRepoConfig(backendRunner)
	backendCommands := git.BackendCommands{
		BackendRunner: backendRunner,
		Config:        &config,
	}
	prodRunner = git.ProdRunner{
		Config:  config,
		Backend: backendCommands,
		Frontend: git.FrontendCommands{
			Frontend: NewFrontendRunner(args.OmitBranchNames, args.DryRun, config.CurrentBranchCache),
			Config:   &config,
			Backend:  &backendCommands,
		},
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
	OmitBranchNames       bool `exhaustruct:"optional"`
	ValidateGitversion    bool `exhaustruct:"optional"`
	ValidateIsRepository  bool `exhaustruct:"optional"`
	ValidateIsConfigured  bool `exhaustruct:"optional"`
	ValidateIsOnline      bool `exhaustruct:"optional"`
}

func NewBackendRunner(dir *string, debug bool) git.BackendRunner {
	backendRunner := subshell.BackendRunner{Dir: dir}
	if debug {
		return subshell.BackendLoggingRunner{Runner: backendRunner}
	}
	return backendRunner
}

// NewFrontendRunner provides a FrontendRunner instance that behaves according to the given configuration.
func NewFrontendRunner(omitBranchNames, dryRun bool, currentBranchCache *cache.String) git.FrontendRunner {
	if dryRun {
		return subshell.FrontendDryRunner{
			CurrentBranch:   currentBranchCache,
			OmitBranchNames: omitBranchNames,
		}
	}
	return subshell.FrontendRunner{
		CurrentBranch:   currentBranchCache,
		OmitBranchNames: omitBranchNames,
	}
}
