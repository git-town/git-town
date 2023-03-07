package cmd

import (
	"fmt"

	"github.com/git-town/git-town/v7/src/cli"
	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
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
			switch args[0] {
			case "add":
				err := addAliases(repo)
				if err != nil {
					cli.Exit(err)
				}
			case "remove":
				err := removeAliases(repo)
				if err != nil {
					cli.Exit(err)
				}
			default:
				cli.Exit(fmt.Errorf(`invalid argument %q. Please provide either "add" or "remove"`, args[0]))
			}
		},
		Args:    cobra.ExactArgs(1),
		GroupID: "setup",
	}
}

func addAliases(repo *git.ProdRepo) error {
	for _, aliasType := range config.AliasTypes() {
		result, err1 := repo.Config.AddGitAlias(aliasType)
		err2 := repo.LoggingShell.PrintCommandAndOutput(result)
		if err1 != nil {
			return err1
		}
		if err2 != nil {
			return err2
		}
	}
	return nil
}

func removeAliases(repo *git.ProdRepo) error {
	for _, aliasType := range config.AliasTypes() {
		existingAlias := repo.Config.GitAlias(aliasType)
		if existingAlias == "town "+string(aliasType) {
			result, err1 := repo.Config.RemoveGitAlias(string(aliasType))
			err2 := repo.LoggingShell.PrintCommandAndOutput(result)
			if err1 != nil {
				return err1
			}
			if err2 != nil {
				return err2
			}
		}
	}
	return nil
}
