package execute

import (
	"github.com/git-town/git-town/v9/src/cache"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/statistics"
	"github.com/git-town/git-town/v9/src/subshell"
	"github.com/git-town/git-town/v9/src/validate"
)

func LoadProdRunner(args LoadArgs) (prodRunner git.ProdRunner, err error) {
	err = validate.HasGitVersion(&prodRunner.Backend)
	if err != nil {
		return
	}
	stats := loadStatistics(args.Debug)
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
	if args.DryRun {
		prodRunner.Config.DryRun = true
	}
	return prodRunner, nil
}

type LoadArgs struct {
	Debug           bool
	DryRun          bool
	OmitBranchNames bool
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

type Statistics interface {
	RegisterRun()
	PrintAnalysis()
}

func loadStatistics(debug bool) Statistics {
	if debug {
		return &statistics.CommandsRun{CommandsCount: 0}
	} else {
		return &statistics.None{}
	}
}
