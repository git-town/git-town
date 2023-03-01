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
	"errors"
	"fmt"

	"github.com/fatih/color"
	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/dialog"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/hosting"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

// Execute runs the Cobra stack.
func Execute() {
	debugFlag := false
	repo := git.NewProdRepo(&debugFlag)
	rootCmd := RootCmd(&repo, &debugFlag)
	majorVersion, minorVersion, err := repo.Silent.Version()
	if err != nil {
		cli.Exit(err)
	}
	if !IsAcceptableGitVersion(majorVersion, minorVersion) {
		cli.Exit(errors.New("this app requires Git 2.7.0 or higher"))
	}
	color.NoColor = false // Prevent color from auto disable
	if err := rootCmd.Execute(); err != nil {
		cli.Exit(err)
	}
}

// RootCmd is the main Cobra object.
func RootCmd(repo *git.ProdRepo, debugFlag *bool) *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "git-town",
		Short: "Generic, high-level Git workflow support",
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
	rootCmd.AddCommand(abortCmd(repo))
	rootCmd.AddCommand(aliasCommand(repo))
	rootCmd.AddCommand(appendCmd(repo))
	rootCmd.AddCommand(completionsCmd(&rootCmd))
	rootCmd.AddCommand(configCmd(repo))
	rootCmd.AddCommand(continueCmd(repo))
	rootCmd.AddCommand(diffParentCommand(repo))
	rootCmd.AddCommand(hackCmd(repo))
	rootCmd.AddCommand(killCommand(repo))
	rootCmd.AddCommand(newPullRequestCommand(repo))
	rootCmd.AddCommand(prependCommand(repo))
	rootCmd.AddCommand(pruneBranchesCommand(repo))
	rootCmd.AddCommand(renameBranchCommand(repo))
	rootCmd.AddCommand(repoCommand(repo))
	rootCmd.AddCommand(statusCommand(repo))
	rootCmd.AddCommand(setParentCommand(repo))
	rootCmd.AddCommand(shipCmd(repo))
	rootCmd.AddCommand(skipCmd(repo))
	rootCmd.AddCommand(switchCmd(repo))
	rootCmd.AddCommand(syncCmd(repo))
	rootCmd.AddCommand(undoCmd(repo))
	rootCmd.AddCommand(versionCmd())
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolVar(debugFlag, "debug", false, "Print all Git commands run under the hood")
	return &rootCmd
}

// IsAcceptableGitVersion indicates whether the given Git version works for Git Town.
func IsAcceptableGitVersion(major, minor int) bool {
	return major > 2 || (major == 2 && minor >= 7)
}

func validateIsConfigured(repo *git.ProdRepo) error {
	err := dialog.EnsureIsConfigured(repo)
	if err != nil {
		return err
	}
	return repo.RemoveOutdatedConfiguration()
}

// ValidateIsRepository asserts that the current directory is in a Git repository.
// If so, it also navigates to the root directory.
func ValidateIsRepository(repo *git.ProdRepo) error {
	if !repo.Silent.IsRepository() {
		return errors.New("this is not a Git repository")
	}
	return repo.NavigateToRootIfNecessary()
}

// handleUnfinishedState checks for unfinished state on disk, handles it, and signals whether to continue execution of the originally intended steps.
//
//nolint:nonamedreturns  // return value isn't obvious from function name
func handleUnfinishedState(repo *git.ProdRepo, connector hosting.Connector) (quit bool, err error) {
	runState, err := runstate.Load(repo)
	if err != nil {
		return false, fmt.Errorf("cannot load previous run state: %w", err)
	}
	if runState == nil || !runState.IsUnfinished() {
		return false, nil
	}
	response, err := dialog.AskHowToHandleUnfinishedRunState(
		runState.Command,
		runState.UnfinishedDetails.EndBranch,
		runState.UnfinishedDetails.EndTime,
		runState.UnfinishedDetails.CanSkip,
	)
	if err != nil {
		return quit, err
	}
	switch response {
	case dialog.ResponseTypeDiscard:
		err = runstate.Delete(repo)
		return false, err
	case dialog.ResponseTypeContinue:
		hasConflicts, err := repo.Silent.HasConflicts()
		if err != nil {
			return false, err
		}
		if hasConflicts {
			return false, fmt.Errorf("you must resolve the conflicts before continuing")
		}
		return true, runstate.Execute(runState, repo, connector)
	case dialog.ResponseTypeAbort:
		abortRunState := runState.CreateAbortRunState()
		return true, runstate.Execute(&abortRunState, repo, connector)
	case dialog.ResponseTypeSkip:
		skipRunState := runState.CreateSkipRunState()
		return true, runstate.Execute(&skipRunState, repo, connector)
	default:
		return false, fmt.Errorf("unknown response: %s", response)
	}
}
