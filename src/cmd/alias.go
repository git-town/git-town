package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/git-town/git-town/src/git"
	"github.com/spf13/cobra"
)

var aliasCommand = &cobra.Command{
	Use:   "alias (true | false)",
	Short: "Adds or removes default global aliases",
	Long: `Adds or removes default global aliases

Global aliases make Git Town commands feel like native Git commands.
When enabled, you can run "git hack" instead of "git town hack".

Does not overwrite existing aliases.

This can conflict with other tools that also define Git aliases.`,
	Run: func(cmd *cobra.Command, args []string) {
		toggle, err := strconv.ParseBool(args[0])
		if err != nil {
			fmt.Printf(`Error: invalid argument: %q. Please provide either "true" or "false".\n`, args[0])
			os.Exit(1)
		}
		var commandsToAlias = []string{
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
		for _, command := range commandsToAlias {
			if toggle {
				addAlias(command, repo())
			} else {
				removeAlias(command, repo())
			}
		}
	},
	Args: cobra.ExactArgs(1),
}

func addAlias(command string, repo *git.ProdRepo) {
	result := repo.AddGitAlias(command)
	repo.LoggingShell.PrintCommand(result.Command(), result.Args()...)
}

func removeAlias(command string, repo *git.ProdRepo) {
	existingAlias := repo.GetGitAlias(command)
	if existingAlias == "town "+command {
		result := repo.RemoveGitAlias(command)
		repo.LoggingShell.PrintCommand(result.Command(), result.Args()...)
	}
}

func init() {
	RootCmd.AddCommand(aliasCommand)
}
