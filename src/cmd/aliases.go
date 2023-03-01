package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/git-town/git-town/v7/src/runstate"
	"github.com/spf13/cobra"
)

func aliasCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:   "aliases (add | remove)",
		Short: "Adds or removes default global aliases",
		Long: `Adds or removes default global aliases

Global aliases make Git Town commands feel like native Git commands.
When enabled, you can run "git hack" instead of "git town hack".

Does not overwrite existing aliases.

This can conflict with other tools that also define Git aliases.`,
		Run: func(cmd *cobra.Command, args []string) {
			var action func(string, *git.ProdRepo) error
			switch args[0] {
			case "add":
				action = addAlias
			case "remove":
				action = removeAlias
			default:
				cli.Exit(fmt.Errorf(`invalid argument %q. Please provide either "add" or "remove"`, args[0]))
			}
			commandsToAlias := []string{
				"append",
				"hack",
				"kill",
				"new-pull-request",
				"prepend",
				"prune-branches",
				"rename-branch",
				"repo",
				"ship",
				"sync",
			}
			ec := runstate.ErrorChecker{}
			for _, command := range commandsToAlias {
				ec.Check(action(command, repo))
			}
			if ec.Err != nil {
				cli.Exit(ec.Err)
			}
		},
		Args:    cobra.ExactArgs(1),
		GroupID: "setup",
	}
}

func addAlias(command string, repo *git.ProdRepo) error {
	result, err := repo.Config.AddGitAlias(command)
	if err != nil {
		return fmt.Errorf("cannot create alias for %q: %w", command, err)
	}
	return repo.LoggingShell.PrintCommand(result.Command(), result.Args()...)
}

func removeAlias(command string, repo *git.ProdRepo) error {
	existingAlias := repo.Config.GitAlias(command)
	if existingAlias == "town "+command {
		result, err := repo.Config.RemoveGitAlias(command)
		if err != nil {
			return fmt.Errorf("cannot remove alias for %q: %w", command, err)
		}
		return repo.LoggingShell.PrintCommand(result.Command(), result.Args()...)
	}
	return nil
}
