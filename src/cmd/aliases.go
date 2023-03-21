package cmd

import (
	"fmt"
	"strings"

	"github.com/git-town/git-town/v7/src/config"
	"github.com/git-town/git-town/v7/src/git"
	"github.com/spf13/cobra"
)

const aliasesSummary = "Adds or removes default global aliases"

const aliasesDesc = `
Global aliases make Git Town commands feel like native Git commands.
When enabled, you can run "git hack" instead of "git town hack".

Does not overwrite existing aliases.

This can conflict with other tools that also define Git aliases.`

func aliasesCommand() *cobra.Command {
	addDebugFlag, readDebugFlag := debugFlag()
	cmd := cobra.Command{
		Use:     "aliases (add | remove)",
		GroupID: "setup",
		Args:    cobra.ExactArgs(1),
		Short:   aliasesSummary,
		Long:    long(aliasesSummary, aliasesDesc),
		RunE: func(cmd *cobra.Command, args []string) error {
			return aliases(args[0], readDebugFlag(cmd))
		},
	}
	addDebugFlag(&cmd)
	return &cmd
}

func aliases(arg string, debug bool) error {
	repo, exit, err := LoadPublicRepo(RepoArgs{
		debug:                 debug,
		dryRun:                false,
		omitBranchNames:       true,
		handleUnfinishedState: false,
		validateGitversion:    true,
	})
	if err != nil || exit {
		return err
	}
	switch strings.ToLower(arg) {
	// TODO: make enum
	case "add":
		return addAliases(&repo)
	case "remove":
		return removeAliases(&repo)
	}
	return fmt.Errorf(`invalid argument %q. Please provide either "add" or "remove"`, arg)
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
