package cmd

import (
	"fmt"

	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/script"
	"github.com/spf13/cobra"
)

var aliasCommand = &cobra.Command{
	Use:   "alias (true | false)",
	Short: "Adds or removes default global aliases",
	Run: func(cmd *cobra.Command, args []string) {
		addOrRemoveAliases(stringToBool(args[0]))
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
	return "!gt " + command
}

func removeAlias(command string) {
	key := getAliasKey(command)
	previousAlias := git.GetGlobalConfigurationValue(key)
	if previousAlias == getAliasValue(command) {
		script.RunCommandSafe("git", "config", "--global", "--unset", key)
	}
}

func addOrRemoveAliases(toggle bool) {
	commands := []string{
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
	for _, command := range commands {
		if toggle {
			addAlias(command)
		} else {
			removeAlias(command)
		}
	}
	fmt.Println() // Trailing newline
}

func init() {
	RootCmd.AddCommand(aliasCommand)
}
