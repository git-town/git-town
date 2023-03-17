package cmd

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

func aliasCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "aliases (add | remove)",
		GroupID: "setup",
		Args:    cobra.ExactArgs(1),
		Short:   "Adds or removes default global aliases",
		Long: `Adds or removes default global aliases

Global aliases make Git Town commands feel like native Git commands.
When enabled, you can run "git hack" instead of "git town hack".

Does not overwrite existing aliases.

This can conflict with other tools that also define Git aliases.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			repo := Repo(debug, false)
			ensure(repo, hasGitVersion),
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

func addAliases(repo *git.PublicRepo) error {
	for _, aliasType := range config.AliasTypes() {
		err := repo.AddGitAlias(aliasType)
		if err != nil {
			return err
		}
	}
	return nil
}

func removeAliases(repo *git.PublicRepo) error {
	for _, aliasType := range config.AliasTypes() {
		existingAlias := repo.Config.GitAlias(aliasType)
		if existingAlias == "town "+string(aliasType) {
			err := repo.RemoveGitAlias(aliasType)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
