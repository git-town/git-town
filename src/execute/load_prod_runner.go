package execute

import (
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
	config := git.NewRepoConfig(backendRunner)
	prodRunner = git.ProdRunner{
		Config: config,
		Backend: git.BackendCommands{
			BackendRunner: backendRunner,
			Config:        &config,
		},
		Frontend: git.FrontendCommands{
			FrontendRunner: NewFrontendRunner(args.OmitBranchNames, args.DryRun, prodRunner.Backend.CurrentBranch, stats),
			Config:         &config,
		},
		Stats: stats,
	}
	err = validate.HasGitVersion(&prodRunner.Backend)
	if err != nil {
		return prodRunner, err
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
