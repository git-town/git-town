package cmd

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v9/src/cli/flags"
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/git"
	"github.com/git-town/git-town/v9/src/messages"
	"github.com/spf13/cobra"
)

const aliasesDesc = "Adds or removes default global aliases"

const aliasesHelp = `
Global aliases make Git Town commands feel like native Git commands.
When enabled, you can run "git hack" instead of "git town hack".

Does not overwrite existing aliases.

This can conflict with other tools that also define Git aliases.`

func aliasesCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := flags.Debug()
	cmd := cobra.Command{
		Use:     "aliases (add | remove)",
		GroupID: "setup",
		Args:    cobra.ExactArgs(1),
		Short:   aliasesDesc,
		Long:    long(aliasesDesc, aliasesHelp),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeAliases(args[0], readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func executeAliases(arg string, debug bool) error {
	repo, err := execute.OpenRepo(execute.OpenRepoArgs{
		Debug:            debug,
		DryRun:           false,
		OmitBranchNames:  true,
		ValidateIsOnline: false,
		ValidateGitRepo:  false,
	})
	if err != nil {
		return err
	}
	switch strings.ToLower(arg) {
	case "add":
		return addAliases(&repo.Runner)
	case "remove":
		return removeAliases(&repo.Runner)
	}
	return fmt.Errorf(messages.InputAddOrRemove, arg)
}

func addAliases(run *git.ProdRunner) error {
	for _, alias := range config.Aliases() {
		err := run.Frontend.AddGitAlias(alias)
		if err != nil {
			return err
		}
	}
	fmt.Printf(messages.CommandsRun, run.CommandsCounter.Count())
	return nil
}

func removeAliases(run *git.ProdRunner) error {
	for _, alias := range config.Aliases() {
		existingAlias := run.Config.GitAlias(alias)
		if existingAlias == "town "+alias.String() {
			err := run.Frontend.RemoveGitAlias(alias)
			if err != nil {
				return err
			}
		}
	}
	fmt.Printf(messages.CommandsRun, run.CommandsCounter.Count())
	return nil
}
