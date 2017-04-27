package cmd

import (
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/script"
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
	Run: func(cmd *cobra.Command, args []string) {
		toggle := stringToBool(args[0])
		for _, command := range commandsToAlias {
			if toggle {
				addAlias(command)
			} else {
				removeAlias(command)
			}
		}
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := validateArgsCount(args, 1)
		if err != nil {
			return err
		}
		return validateBooleanArgument(args[0])
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
