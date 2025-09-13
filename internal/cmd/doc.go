// Package cmd defines the Git Town commands.
package cmd

import (
	"context"
	"time"

	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/cmd/config"
	"github.com/git-town/git-town/v21/internal/cmd/ship"
	"github.com/git-town/git-town/v21/internal/cmd/status"
	"github.com/git-town/git-town/v21/internal/cmd/swap"
	"github.com/git-town/git-town/v21/internal/cmd/sync"
	"github.com/git-town/git-town/v21/internal/config/cliconfig"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/execute"
	"github.com/git-town/git-town/v21/internal/update"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Execute runs the Cobra stack.
func Execute() error {
	rootCmd := rootCmd()
	rootCmd.AddCommand(appendCmd())
	rootCmd.AddCommand(branchCmd())
	rootCmd.AddCommand(completionsCmd(&rootCmd))
	rootCmd.AddCommand(compressCmd())
	rootCmd.AddCommand(config.RootCmd())
	rootCmd.AddCommand(continueCmd())
	rootCmd.AddCommand(contributeCmd())
	rootCmd.AddCommand(diffParentCommand())
	rootCmd.AddCommand(hackCmd())
	rootCmd.AddCommand(detachCommand())
	rootCmd.AddCommand(deleteCommand())
	rootCmd.AddCommand(downCmd())
	rootCmd.AddCommand(featureCmd())
	rootCmd.AddCommand(initCommand())
	rootCmd.AddCommand(mergeCommand())
	rootCmd.AddCommand(observeCmd())
	rootCmd.AddCommand(offlineCmd())
	rootCmd.AddCommand(parkCmd())
	rootCmd.AddCommand(proposeCommand())
	rootCmd.AddCommand(prependCommand())
	rootCmd.AddCommand(prototypeCmd())
	rootCmd.AddCommand(renameCommand())
	rootCmd.AddCommand(repoCommand())
	rootCmd.AddCommand(runLogCommand())
	rootCmd.AddCommand(status.RootCommand())
	rootCmd.AddCommand(setParentCommand())
	rootCmd.AddCommand(ship.Cmd())
	rootCmd.AddCommand(skipCmd())
	rootCmd.AddCommand(swap.Cmd())
	rootCmd.AddCommand(switchCmd())
	rootCmd.AddCommand(sync.Cmd())
	rootCmd.AddCommand(undoCmd())
	rootCmd.AddCommand(upCmd())
	rootCmd.AddCommand(walkCommand())
	err := rootCmd.Execute()

	checkForUpdates()

	return err
}

// checkForUpdates checks for available updates and displays a notification
func checkForUpdates() {
	// Create a timeout context for the update check
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Try to open the repo to get configuration, but don't fail if it's not a repo
	cliConfig := cliconfig.New(cliconfig.NewArgs{
		AutoResolve: None[configdomain.AutoResolve](),
		Detached:    None[configdomain.Detached](),
		DryRun:      None[configdomain.DryRun](),
		Stash:       None[configdomain.Stash](),
		Verbose:     None[configdomain.Verbose](),
	})

	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		CliConfig:        cliConfig,
		PrintBranchNames: false,
		PrintCommands:    false,
		ValidateGitRepo:  false,
		ValidateIsOnline: false,
	})
	if err != nil {
		// If we can't open the repo, skip update check
		return
	}

	// Check if update checks are enabled
	if repo.UnvalidatedConfig.NormalConfig.UpdateCheck.IsDisabled() {
		return
	}

	// Skip update check if offline
	if repo.UnvalidatedConfig.NormalConfig.Offline.IsOffline() {
		return
	}

	// Create notifier and check for updates
	logger := print.Logger{}
	notifier := update.NewNotifier(logger)
	notifier.CheckAndNotify(ctx)
}
