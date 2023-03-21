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
	"fmt"

	"github.com/git-town/git-town/v7/src/cache"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/git-town/git-town/v7/src/subshell"
	"github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

func rootCmd() cobra.Command {
	rootCmd := cobra.Command{
		Use:           "git-town",
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "Generic, high-level Git workflow support",
		Long: `Git Town makes software development teams who use Git even more productive and happy.

It adds Git commands that support GitHub Flow, Git Flow, the Nvie model, GitLab Flow, and other workflows more directly,
and it allows you to perform many common Git operations faster and easier.`,
	}
	rootCmd.AddGroup(&cobra.Group{
		ID:    "basic",
		Title: "Basic commands:",
	}, &cobra.Group{
		ID:    "errors",
		Title: "Commands to deal with errors:",
	}, &cobra.Group{
		ID:    "lineage",
		Title: "Commands for nested feature branches:",
	}, &cobra.Group{
		ID:    "setup",
		Title: "Commands to set up Git Town on your computer:",
	})
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	return rootCmd
}

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

// boolFlag provides access to boolean Cobra command-line flags
// in a way where Go's usage checker (which produces compilation errors for unused variables)
// enforces that the flag is guaranteed to be defined and used.
// This reduces programmer errors while defining and using command-line flags..
func boolFlag(name, short string, desc string) (func(*cobra.Command), func(*cobra.Command) bool) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().BoolP(name, short, false, desc)
	}
	readFlag := func(cmd *cobra.Command) bool {
		value, err := cmd.Flags().GetBool(name)
		if err != nil {
			panic(fmt.Sprintf("command %q does not have a boolean %q flag", cmd.Name(), name))
		}
		return value
	}
	return addFlag, readFlag
}

// stringFlag provides access to boolean Cobra command-line flags
// in a way where Go's usage checker (which produces compilation errors for unused variables)
// enforces that the flag is guaranteed to be defined and used.
// This reduces programmer errors while defining and using command-line flags..
func stringFlag(name, short, defaultValue, desc string) (func(*cobra.Command), func(*cobra.Command) string) {
	addFlag := func(cmd *cobra.Command) {
		cmd.PersistentFlags().StringP(name, short, defaultValue, desc)
	}
	readFlag := func(cmd *cobra.Command) string {
		value, err := cmd.Flags().GetString(name)
		if err != nil {
			panic(fmt.Sprintf("command %q does not have a string %q flag", cmd.Name(), name))
		}
		return value
	}
	return addFlag, readFlag
}

func debugFlag() (func(*cobra.Command), func(*cobra.Command) bool) {
	return boolFlag("debug", "d", "Print all Git commands run under the hood")
}

func dryRunFlag() (func(*cobra.Command), func(*cobra.Command) bool) {
	return boolFlag("dry-run", "", "Print but do not run the Git commands")
}

func LoadPublicRepo(args RepoArgs) (repo git.PublicRepo, exit bool, err error) { //nolint:nonamedreturns // so many return values require names
	internalRepo := internalRepo(args.debug)
	repo = publicRepo(args.omitBranchNames, args.dryRun, &internalRepo)
	if args.validateIsRepository {
		err := validate.IsRepository(&repo)
		if err != nil {
			return repo, false, err
		}
	}
	if !args.omitBranchNames || args.dryRun {
		currentBranch, err := internalRepo.CurrentBranch()
		if err != nil {
			return repo, false, err
		}
		internalRepo.CurrentBranchCache.Set(currentBranch)
	}
	if args.dryRun {
		internalRepo.DryRun = true
	}
	ec := runstate.ErrorChecker{}
	if args.validateGitversion {
		ec.Check(validate.HasGitVersion(&internalRepo))
	}
	if args.validateIsConfigured {
		ec.Check(validate.IsConfigured(&repo))
	}
	if args.validateIsOnline {
		ec.Check(validate.IsOnline(&repo))
	}
	if args.handleUnfinishedState {
		exit = ec.Bool(validate.HandleUnfinishedState(&repo, nil))
	}
	return repo, exit, ec.Err
}

type RepoArgs struct {
	debug                 bool
	dryRun                bool
	handleUnfinishedState bool
	omitBranchNames       bool `exhaustruct:"optional"`
	validateGitversion    bool `exhaustruct:"optional"`
	validateIsRepository  bool `exhaustruct:"optional"`
	validateIsConfigured  bool `exhaustruct:"optional"`
	validateIsOnline      bool `exhaustruct:"optional"`
}

func internalRepo(debug bool) git.InternalRepo {
	shellRunner := subshell.InternalRunner{Dir: nil}
	var gitRunner git.InternalRunner
	if debug {
		gitRunner = subshell.InternalDebuggingRunner{InternalRunner: shellRunner}
	} else {
		gitRunner = shellRunner
	}
	return git.InternalRepo{
		InternalRunner:     gitRunner,
		Config:             config.NewGitTown(gitRunner),
		CurrentBranchCache: &cache.String{},
		DryRun:             false,
		IsRepoCache:        &cache.Bool{},
		RemoteBranchCache:  &cache.Strings{},
		RemotesCache:       &cache.Strings{},
		RootDirCache:       &cache.String{},
	}
}

func publicRepo(omitBranchNames, dryRun bool, internalRepo *git.InternalRepo) git.PublicRepo {
	var gitRunner git.PublicRunner
	if dryRun {
		gitRunner = subshell.PublicDryRunner{
			CurrentBranch:   internalRepo.CurrentBranchCache,
			OmitBranchNames: omitBranchNames,
		}
	} else {
		gitRunner = subshell.PublicRunner{
			CurrentBranch:   internalRepo.CurrentBranchCache,
			OmitBranchNames: omitBranchNames,
		}
	}
	return git.PublicRepo{
		Public:       gitRunner,
		InternalRepo: *internalRepo,
	}
}
