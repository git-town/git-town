// Package cmd defines the Git Town commands.
package cmd

import (
	"github.com/git-town/git-town/v16/internal/cmd/config"
	"github.com/git-town/git-town/v16/internal/cmd/debug"
	"github.com/git-town/git-town/v16/internal/cmd/ship"
	"github.com/git-town/git-town/v16/internal/cmd/status"
	"github.com/git-town/git-town/v16/internal/cmd/sync"
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
	rootCmd.AddCommand(debug.RootCmd())
	rootCmd.AddCommand(diffParentCommand())
	rootCmd.AddCommand(hackCmd())
	rootCmd.AddCommand(deleteCommand())
	rootCmd.AddCommand(killCommand())
	rootCmd.AddCommand(newPullRequestCommand())
	rootCmd.AddCommand(observeCmd())
	rootCmd.AddCommand(offlineCmd())
	rootCmd.AddCommand(parkCmd())
	rootCmd.AddCommand(proposeCommand())
	rootCmd.AddCommand(prependCommand())
	rootCmd.AddCommand(prototypeCmd())
	rootCmd.AddCommand(renameCommand())
	rootCmd.AddCommand(repoCommand())
	rootCmd.AddCommand(status.RootCommand())
	rootCmd.AddCommand(setParentCommand())
	rootCmd.AddCommand(ship.Cmd())
	rootCmd.AddCommand(skipCmd())
	rootCmd.AddCommand(switchCmd())
	rootCmd.AddCommand(sync.Cmd())
	rootCmd.AddCommand(undoCmd())
	return rootCmd.Execute()
}
