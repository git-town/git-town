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
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/subshell"
	"github.com/git-town/git-town/v7/src/validate"
	"github.com/spf13/cobra"
)

// Execute runs the Cobra stack.
func Execute() error {
	rootCmd := rootCmd()

	rootCmd.AddCommand(abortCmd())
	rootCmd.AddCommand(aliasCommand())
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

func debugFlag(cmd *cobra.Command, flag *bool) {
	cmd.PersistentFlags().BoolVar(flag, "debug", false, "Print all Git commands run under the hood")
}

func dryRunFlag(cmd *cobra.Command) *bool {
	return cmd.PersistentFlags().Bool("dryrun", false, "Print but do not execute the Git commands")
}

func Repo(debug, dryRun bool) git.PublicRepo {
	internalRepo := internalRepo(debug)
	return publicRepo(dryRun, &internalRepo)
}

func internalRepo(debug bool) git.InternalRepo {
	shellRunner := subshell.InternalRunner{}
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
		DryRun:             &subshell.DryRun{},
		IsRepoCache:        &cache.Bool{},
		RemoteBranchCache:  &cache.Strings{},
		RemotesCache:       &cache.Strings{},
		RootDirCache:       &cache.String{},
	}
}

func publicRepo(dryRun bool, internalRepo *git.InternalRepo) git.PublicRepo {
	var gitRunner git.PublicRunner
	if dryRun {
		gitRunner = subshell.PublicDryRunner{}
	} else {
		gitRunner = subshell.PublicRunner{}
	}
	return git.PublicRepo{
		Public:       gitRunner,
		InternalRepo: *internalRepo,
	}
}

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

var (
	ensure        = validate.Ensure        //nolint:gochecknoglobals
	hasGitVersion = validate.HasGitVersion //nolint:gochecknoglobals
	isRepository  = validate.IsRepository  //nolint:gochecknoglobals
	isConfigured  = validate.IsConfigured  //nolint:gochecknoglobals
	isOnline      = validate.IsOnline      //nolint:gochecknoglobals
)
