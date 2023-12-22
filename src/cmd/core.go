// Package cmd defines the Git Town commands.
package cmd

import (
	"github.com/git-town/git-town/v11/src/cmd/configcmds"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/vm/opcode"
	"github.com/git-town/git-town/v11/src/vm/program"
)

// Execute runs the Cobra stack.
func Execute() error {
	rootCmd := rootCmd()
	rootCmd.AddCommand(aliasesCommand())
	rootCmd.AddCommand(appendCmd())
	rootCmd.AddCommand(completionsCmd(&rootCmd))
	rootCmd.AddCommand(configcmds.ConfigCmd())
	rootCmd.AddCommand(continueCmd())
	rootCmd.AddCommand(diffParentCommand())
	rootCmd.AddCommand(hackCmd())
	rootCmd.AddCommand(killCommand())
	rootCmd.AddCommand(newPullRequestCommand())
	rootCmd.AddCommand(proposeCommand())
	rootCmd.AddCommand(prependCommand())
	rootCmd.AddCommand(renameBranchCommand())
	rootCmd.AddCommand(repoCommand())
	rootCmd.AddCommand(statusCommand())
	rootCmd.AddCommand(setParentCommand())
	rootCmd.AddCommand(shipCmd())
	rootCmd.AddCommand(skipCmd())
	rootCmd.AddCommand(switchCmd())
	rootCmd.AddCommand(syncCmd())
	rootCmd.AddCommand(undoCmd())
	return rootCmd.Execute()
}

// wrap wraps the given list with opcodes that change the Git root directory or stash away open changes.
// TODO: only wrap if the list actually contains any opcodes.
func wrap(program *program.Program, options wrapOptions) {
	program.Add(&opcode.PreserveCheckoutHistory{
		PreviousBranchCandidates: options.PreviousBranchCandidates,
	})
	if options.StashOpenChanges {
		program.Prepend(&opcode.StashOpenChanges{})
		program.Add(&opcode.RestoreOpenChanges{})
	}
}

// wrapOptions represents the options given to Wrap.
type wrapOptions struct {
	RunInGitRoot             bool
	StashOpenChanges         bool
	PreviousBranchCandidates gitdomain.LocalBranchNames
}
