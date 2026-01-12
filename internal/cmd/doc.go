// Package cmd defines the Git Town commands.
package cmd

import (
	"github.com/git-town/git-town/v22/internal/cmd/config"
	"github.com/git-town/git-town/v22/internal/cmd/ship"
	"github.com/git-town/git-town/v22/internal/cmd/status"
	"github.com/git-town/git-town/v22/internal/cmd/swap"
	"github.com/git-town/git-town/v22/internal/cmd/sync"
)

// Execute runs the Cobra stack.
func Execute() error {
	rootCmd := rootCmd()
	rootCmd.AddCommand(appendCmd())
	rootCmd.AddCommand(branchCmd())
	rootCmd.AddCommand(commitCmd())
	rootCmd.AddCommand(completionsCmd(&rootCmd))
	rootCmd.AddCommand(compressCmd())
	rootCmd.AddCommand(config.RootCmd())
	rootCmd.AddCommand(continueCmd())
	rootCmd.AddCommand(contributeCmd())
	rootCmd.AddCommand(diffParentCommand())
	rootCmd.AddCommand(detachCommand())
	rootCmd.AddCommand(deleteCommand())
	rootCmd.AddCommand(downCmd())
	rootCmd.AddCommand(featureCmd())
	rootCmd.AddCommand(hackCmd())
	rootCmd.AddCommand(initCommand())
	rootCmd.AddCommand(mergeCommand())
	rootCmd.AddCommand(observeCmd())
	rootCmd.AddCommand(offlineCmd())
	rootCmd.AddCommand(parkCmd())
	rootCmd.AddCommand(prependCommand())
	rootCmd.AddCommand(proposeCommand())
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
	return rootCmd.Execute()
}
