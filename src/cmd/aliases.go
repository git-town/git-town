package cmd

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v11/src/cli/flags"
	"github.com/git-town/git-town/v11/src/cmd/cmdhelpers"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/execute"
	"github.com/git-town/git-town/v11/src/git"
	"github.com/git-town/git-town/v11/src/messages"
	"github.com/spf13/cobra"
)

const aliasesDesc = "Adds or removes default global aliases"

const aliasesHelp = `
Global aliases make Git Town commands feel like native Git commands. When enabled, you can run "git hack" instead of "git town hack".

Does not overwrite existing aliases.

This can conflict with other tools that also define Git aliases.`

func aliasesCommand() *cobra.Command {
	addVerboseFlag, readVerboseFlag := flags.Verbose()
	addDryRunFlag, readDryRunFlag := flags.DryRun()
	cmd := cobra.Command{
		Use:     "aliases (add | remove)",
		GroupID: "setup",
		Args:    cobra.ExactArgs(1),
		Short:   aliasesDesc,
		Long:    cmdhelpers.Long(aliasesDesc, aliasesHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeAliases(args[0], readDryRunFlag(cmd), readVerboseFlag(cmd))
		},
	}
	addDryRunFlag(&cmd)
	addVerboseFlag(&cmd)
	return &cmd
}

func executeAliases(arg string, dryRun, verbose bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Verbose:          verbose,
		DryRun:           dryRun,
		OmitBranchNames:  true,
		PrintCommands:    true,
		ValidateIsOnline: false,
		ValidateGitRepo:  false,
	})
	if err != nil {
		return err
	}
	switch strings.ToLower(arg) {
	case "add":
		return addAliases(repo.Runner)
	case "remove":
		return removeAliases(repo.Runner)
	}
	return fmt.Errorf(messages.InputAddOrRemove, arg)
}

func addAliases(run *git.ProdRunner) error {
	for _, aliasableCommand := range configdomain.AllAliasableCommands() {
		existingAlias := run.Aliases[aliasableCommand]
		if existingAlias != "" {
			continue
		}
		err := run.Frontend.SetGitAlias(aliasableCommand)
		if err != nil {
			return err
		}
	}
	fmt.Printf(messages.CommandsRun, run.CommandsCounter.Count())
	return nil
}

func removeAliases(run *git.ProdRunner) error {
	for _, aliasableCommand := range configdomain.AllAliasableCommands() {
		existingAlias := run.Aliases[aliasableCommand]
		if existingAlias == "town "+aliasableCommand.String() {
			err := run.Frontend.RemoveGitAlias(aliasableCommand)
			if err != nil {
				return err
			}
		}
	}
	fmt.Printf(messages.CommandsRun, run.CommandsCounter.Count())
	return nil
}
