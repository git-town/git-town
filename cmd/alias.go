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
		addOrRemoveAliases(convertStringToBool(args[0]))
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := validateArgsCount(args, 1)
		if err != nil {
			return err
		}
		return validateBooleanArgument(args[0])
	},
}

func addOrRemoveAlias(command string, toggle bool) {
	key := "alias." + command
	value := "!gt " + command
	if toggle {
		script.RunCommandSafe("git", "config", "--global", key, value)
	} else {
		previousAlias := git.GetGlobalConfigurationValue(key)
		if previousAlias == value {
			script.RunCommandSafe("git", "config", "--global", "--unset", key)
		}
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
		addOrRemoveAlias(command, toggle)
	}
	fmt.Println() // Trailing newline
}

func init() {
	RootCmd.AddCommand(aliasCommand)
}
