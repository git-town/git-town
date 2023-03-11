package cmd

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func aliasCommand(repo *git.ProdRepo) *cobra.Command {
	return &cobra.Command{
		Use:     "aliases (add | remove)",
		GroupID: "setup",
		Args:    cobra.ExactArgs(1),
		PreRunE: Ensure(repo, HasGitVersion),
		Short:   "Adds or removes default global aliases",
		Long: `Adds or removes default global aliases

Global aliases make Git Town commands feel like native Git commands.
When enabled, you can run "git hack" instead of "git town hack".

Does not overwrite existing aliases.

This can conflict with other tools that also define Git aliases.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			switch strings.ToLower(args[0]) {
			case "add":
				return addAliases(repo)
			case "remove":
				return removeAliases(repo)
			}
			return fmt.Errorf(`invalid argument %q. Please provide either "add" or "remove"`, args[0])
		},
	}
}

func addAliases(repo *git.ProdRepo) error {
	for _, aliasType := range config.AliasTypes() {
		result, err1 := repo.Config.AddGitAlias(aliasType)
		err2 := repo.LoggingRunner.PrintCommandAndOutput(result)
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
			err2 := repo.LoggingRunner.PrintCommandAndOutput(result)
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
