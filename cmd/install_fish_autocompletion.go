package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"

	"github.com/spf13/cobra"
)

var installFishAutocompletionCommand = &cobra.Command{
	Use:   "install-fish-autocompletion",
	Short: "Installs the autocompletion definition for Fish shell (http://fishshell.com)",
	Run: func(cmd *cobra.Command, args []string) {
		installFishAutocompletion()
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateMaxArgs(args, 0)
	},
}

func installFishAutocompletion() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	filename := path.Join(user.HomeDir, ".config", "fish", "completions", "git.fish")
	err = os.MkdirAll(path.Dir(filename), 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filename, []byte(buildAutocompletionDefinition()), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

var fishAutocompletionTemplate = `
# All Git Town commands\n
set git_town_commands %s

# Indicates through its error code whether the command line to auto-complete
# already contains a Git Town command or not.
#
# - doesn't have command yet: exit code 0
# - has command already: exit code 1
function __fish_complete_git_town_no_command
  for cmd in (commandline -otc)
    if contains $cmd $git_town_commands
      return 1
    end
  end
  return 0
end


# Define autocompletion for the Git Town commands themselves.
#
# These only get autocompleted if there is no Git Town command present in the
# command line already.
# This is done through __fish_complete_git_town_no_command
%s


# Define autocompletion of Git branch names.
#
# This is only enabled for commands that take branch names.
# This is achieved through __fish_complete_git_town_command_takes_branch
complete --command git --arguments "(git branch | tr -d '* ')" --no-files


# Define autocompletion for command-line switches
%s
`

type commandDefinition struct {
	name        string
	description string
}

type optionDefinition struct {
	name        string
	description string
}

func buildAutocompletionDefinition() string {
	commands := []commandDefinition{
		commandDefinition{name: "hack", description: hackCmd.Short},
		commandDefinition{name: "kill", description: killCommand.Short},
		commandDefinition{name: "new-pull-request", description: newPullRequestCommand.Short},
		commandDefinition{name: "prune-branches", description: pruneBranchesCommand.Short},
		commandDefinition{name: "rename-branch", description: renameBranchCommand.Short},
		commandDefinition{name: "repo", description: repoCommand.Short},
		commandDefinition{name: "ship", description: shipCmd.Short},
		commandDefinition{name: "sync", description: syncCmd.Short},
	}
	options := []optionDefinition{
		optionDefinition{name: "abort", description: abortFlagDescription},
		optionDefinition{name: "continue", description: continueFlagDescription},
		optionDefinition{name: "undo", description: undoFlagDescription},
	}

	commandsSpaceSeparated := ""
	for _, command := range commands {
		commandsSpaceSeparated += command.name + " "
	}
	commandAutcompletion := ""
	for _, command := range commands {
		commandAutcompletion += fmt.Sprintf("complete --command git --arguments '%s' --description '%s' --condition '__fish_complete_git_town_no_command' --no-files\n", command.name, command.description)
	}
	optionAutocompletion := ""
	for _, option := range options {
		optionAutocompletion += fmt.Sprintf("complete --command git --long-option '%s' --description '%s' --no-files\n", option.name, option.description)
	}

	return fmt.Sprintf(fishAutocompletionTemplate, commandsSpaceSeparated, commandAutcompletion, optionAutocompletion)
}

func init() {
	RootCmd.AddCommand(installFishAutocompletionCommand)
}
