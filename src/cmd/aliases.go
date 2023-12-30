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
Global aliases make Git Town commands feel like native Git commands.
When enabled, you can run "git hack" instead of "git town hack".

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
	fc := execute.FailureCollector{}
	fc.Check(addAlias(configdomain.KeyAliasAppend, run.AliasAppend, "town append", &run.Frontend))
	fc.Check(addAlias(configdomain.KeyAliasDiffParent, run.AliasDiffParent, "town diff-parent", &run.Frontend))
	fc.Check(addAlias(configdomain.KeyAliasHack, run.AliasHack, "town hack", &run.Frontend))
	fc.Check(addAlias(configdomain.KeyAliasKill, run.AliasKill, "town kill", &run.Frontend))
	fc.Check(addAlias(configdomain.KeyAliasPrepend, run.AliasPrepend, "town prepend", &run.Frontend))
	fc.Check(addAlias(configdomain.KeyAliasPropose, run.AliasPropose, "town propose", &run.Frontend))
	fc.Check(addAlias(configdomain.KeyAliasRenameBranch, run.AliasRenameBranch, "town rename-branch", &run.Frontend))
	fc.Check(addAlias(configdomain.KeyAliasRepo, run.AliasRepo, "town repo", &run.Frontend))
	fc.Check(addAlias(configdomain.KeyAliasShip, run.AliasShip, "town ship", &run.Frontend))
	fc.Check(addAlias(configdomain.KeyAliasSync, run.AliasSync, "town sync", &run.Frontend))
	if fc.Err == nil {
		fmt.Printf(messages.CommandsRun, run.CommandsCounter.Count())
	}
	return fc.Err
}

func addAlias(key configdomain.Key, existingValue, newValue string, frontend *git.FrontendCommands) error {
	if existingValue != "" {
		return nil
	}
	return frontend.SetGitAlias(key, newValue)
}

func removeAliases(run *git.ProdRunner) error {
	fc := execute.FailureCollector{}
	fc.Check(removeAlias(configdomain.KeyAliasAppend, run.AliasAppend, "town append", &run.Frontend))
	if fc.Err == nil {
		fmt.Printf(messages.CommandsRun, run.CommandsCounter.Count())
	}
	return fc.Err
}

func removeAlias(key configdomain.Key, existingValue, expectedValue string, frontend *git.FrontendCommands) error {
	if existingValue != expectedValue {
		return nil
	}
	return frontend.RemoveGitAlias(key)
}
