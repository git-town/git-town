package cmd

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/execute"
	"github.com/git-town/git-town/v9/src/flags"
	"github.com/git-town/git-town/v9/src/git"
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
			return aliases(args[0], readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func aliases(arg string, debug bool) error {
	run, exit, err := execute.LoadProdRunner(execute.LoadArgs{
		Debug:                 debug,
		DryRun:                false,
		OmitBranchNames:       true,
		HandleUnfinishedState: false,
		ValidateGitversion:    true,
	})
	if err != nil || exit {
		return err
	}
	switch strings.ToLower(arg) {
	case "add":
		return addAliases(&run)
	case "remove":
		return removeAliases(&run)
	}
	return fmt.Errorf(`invalid argument %q. Please provide either "add" or "remove"`, arg)
}

func addAliases(run *git.ProdRunner) error {
	for _, aliasType := range config.AliasTypes() {
		err := run.Frontend.AddGitAlias(aliasType)
		if err != nil {
			return err
		}
	}
	run.Stats.PrintAnalysis()
	return nil
}

func removeAliases(run *git.ProdRunner) error {
	for _, aliasType := range config.AliasTypes() {
		existingAlias := run.Config.GitAlias(aliasType)
		if existingAlias == "town "+string(aliasType) {
			err := run.Frontend.RemoveGitAlias(aliasType)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
