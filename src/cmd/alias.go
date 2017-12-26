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
	script.RunCommandSafe("git", "config", "--global", getAliasKey(command), getAliasValue(command))
}

func getAliasKey(command string) string {
	return "alias." + command
}

func getAliasValue(command string) string {
	return "town " + command
}

func removeAlias(command string) {
	key := getAliasKey(command)
	previousAlias := git.GetGlobalConfigurationValue(key)
	if previousAlias == getAliasValue(command) {
		script.RunCommandSafe("git", "config", "--global", "--unset", key)
	}
}

func init() {
	RootCmd.AddCommand(aliasCommand)
}
