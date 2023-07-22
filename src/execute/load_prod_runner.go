package execute

import (
	"github.com/git-town/git-town/v9/src/cache"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/failure"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/statistics"
	"github.com/git-town/git-town/v9/src/subshell"
	"github.com/git-town/git-town/v9/src/validate"
)

type Statistics interface {
	RegisterRun()
	PrintAnalysis()
}

func LoadProdRunner(args LoadArgs) (prodRunner git.ProdRunner, exit bool, err error) {
	var stats Statistics
	if args.Debug {
		stats = &statistics.CommandsRun{CommandsCount: 0}
	} else {
		stats = &statistics.None{}
	}
	backendRunner := subshell.BackendRunner{Dir: nil, Verbose: args.Debug, Stats: stats}
	config := git.RepoConfig{
		GitTown: config.NewGitTown(backendRunner),
		DryRun:  false, // to bootstrap this, DryRun always gets initialized as false and later enabled if needed
	}
	backendCommands := git.BackendCommands{
		BackendRunner:      backendRunner,
		Config:             &config,
		CurrentBranchCache: &cache.String{},
		IsRepoCache:        &cache.Bool{},
		RemoteBranchCache:  &cache.Strings{},
		RemotesCache:       &cache.Strings{},
		RootDirCache:       &cache.String{},
	}
	prodRunner = git.ProdRunner{
		Config:  config,
		Backend: backendCommands,
		Frontend: git.FrontendCommands{
			FrontendRunner:         NewFrontendRunner(args.OmitBranchNames, args.DryRun, prodRunner.Backend.CurrentBranch, stats),
			SetCachedCurrentBranch: backendCommands.CurrentBranchCache.Set,
		},
		Stats: stats,
	}
	if args.ValidateIsRepository {
		err := validate.IsRepository(&prodRunner)
		if err != nil {
			return prodRunner, false, err
		}
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
func NewFrontendRunner(omitBranchNames, dryRun bool, getCurrentBranch subshell.GetCurrentBranchFunc, stats Statistics) git.FrontendRunner {
	if dryRun {
		return &subshell.FrontendDryRunner{
			GetCurrentBranch: getCurrentBranch,
			OmitBranchNames:  omitBranchNames,
			Stats:            stats,
		}
	}
	return &subshell.FrontendRunner{
		GetCurrentBranch: getCurrentBranch,
		OmitBranchNames:  omitBranchNames,
		Stats:            stats,
	}
}
