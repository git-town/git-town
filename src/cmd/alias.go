package cmd

import (
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/script"
	"github.com/Originate/git-town/src/util"
	"github.com/spf13/cobra"
)

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

var aliasCommand = &cobra.Command{
	Use:   "alias (true | false)",
	Short: "Adds or removes default global aliases",
	Long: `Adds or removes default global aliases

Global aliases allow Git Town commands to be used like native Git commands.
When aliases are set, you can run "git hack" instead of having to run "git town hack".
Example: "git append" becomes equivalent to "git town append".

When adding aliases, no existing aliases will be overwritten.

Note that this can conflict with other tools that also define additional Git commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		toggle := util.StringToBool(args[0])
		for _, command := range commandsToAlias {
			if toggle {
				addAlias(command)
			} else {
				removeAlias(command)
			}
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			return validateBooleanArgument(args[0])
		}
		return cobra.ExactArgs(1)(cmd, args)
	},
}

func addAlias(command string) {
	result := git.Config().AddGitAlias(command)
	script.PrintCommand(result.Command(), result.Args()...)
}

func removeAlias(command string) {
	existingAlias := git.Config().GetGitAlias(command)
	if existingAlias == "town "+command {
		result := git.Config().RemoveGitAlias(command)
		script.PrintCommand(result.Command(), result.Args()...)
	}
}

func init() {
	RootCmd.AddCommand(aliasCommand)
}
