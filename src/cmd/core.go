// Package cmd defines the Git Town commands.
//
// Each Git Town command begins by inspecting the current state of the Git
// repository (which branch you are on, whether you have open changes). If there
// are no errors, it generates a StepList instance containing the steps to run.
//
// Steps, located in src/steps, implement the individual steps that
// each Git Town command performs. Examples are steps to
// change to a different Git branch or to pull updates for the current branch.
//
// When executing a step, the runstate.Execute function goes through each step in the StepList.
// It executes the step. If it succeeded, it asks the current step to provide the undo step
// for what it just did and appends it to the undo StepList.
// If a Git command fails (for example due to a merge conflict), then the program
// asks the step to create it's corresponding abort and continue steps, adds them to the respective StepLists,
// saves the entire runstate to disk, informs the user, and exits.
//
// When running "git town continue", Git Town loads the runstate and executes the "continue" StepList in it.
// When running "git town abort", Git Town loads the runstate and executes the "abort" StepList in it.
// When running "git town undo", Git Town loads the runstate and executes the "undo" StepList in it.
package cmd

import (
	"github.com/git-town/git-town/v7/src/cache"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/subshell"
	"github.com/git-town/git-town/v7/src/validate"
)

// Execute runs the Cobra stack.
func Execute() error {
	rootCmd := rootCmd()
	rootCmd.AddCommand(abortCmd())
	rootCmd.AddCommand(aliasesCommand())
	rootCmd.AddCommand(appendCmd())
	rootCmd.AddCommand(completionsCmd(&rootCmd))
	rootCmd.AddCommand(configCmd())
	rootCmd.AddCommand(continueCmd())
	rootCmd.AddCommand(diffParentCommand())
	rootCmd.AddCommand(hackCmd())
	rootCmd.AddCommand(killCommand())
	rootCmd.AddCommand(newPullRequestCommand())
	rootCmd.AddCommand(prependCommand())
	rootCmd.AddCommand(pruneBranchesCommand())
	rootCmd.AddCommand(renameBranchCommand())
	rootCmd.AddCommand(repoCommand())
	rootCmd.AddCommand(statusCommand())
	rootCmd.AddCommand(setParentCommand())
	rootCmd.AddCommand(shipCmd())
	rootCmd.AddCommand(skipCmd())
	rootCmd.AddCommand(switchCmd())
	rootCmd.AddCommand(syncCmd())
	rootCmd.AddCommand(undoCmd())
	rootCmd.AddCommand(versionCmd())
	return rootCmd.Execute()
}

func long(summary string, desc ...string) string {
	if len(desc) == 1 {
		return summary + ".\n" + desc[0]
	}
	return summary + "."
}

func LoadProdRunner(args loadArgs) (prodRunner git.ProdRunner, exit bool, err error) { //nolint:nonamedreturns // so many return values require names
	var stats *subshell.Statistics
	if args.debug {
		stats = &subshell.Statistics{}
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

func NewBackendRunner(dir *string, debug bool, statistics *subshell.Statistics) git.BackendRunner {
	backendRunner := subshell.BackendRunner{Dir: dir, Statistics: statistics}
	if debug {
		return subshell.BackendLoggingRunner{Runner: backendRunner, Statistics: statistics}
	}
	return backendRunner
}

// NewFrontendRunner provides a FrontendRunner instance that behaves according to the given configuration.
func NewFrontendRunner(omitBranchNames, dryRun bool, currentBranchCache *cache.String, stats *subshell.Statistics) git.FrontendRunner {
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
