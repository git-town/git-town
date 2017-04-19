package prompt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Originate/git-town/lib/git"
	"github.com/fatih/color"
)

func EnsureKnowsConfiguration() {
	printConfigurationHeader()

	newMainBranch := askForBranch(branchPromptConfig{
		branchNames: git.GetLocalBranches(),
		prompt:      getMainBranchPrompt(),
		validate: func(branchName string) error {
			if branchName == "" {
				return errors.New("A main development branch is required to enable the features provided by Git Town")
			}
			return nil
		},
	})
	git.SetMainBranch(newMainBranch)

	var newPerennialBranches []string
	for {
		newPerennialBranch := askForBranch(branchPromptConfig{
			branchNames: git.GetLocalBranches(),
			prompt:      getPerennialBranchesPrompt(),
			validate: func(branchName string) error {
				if branchName == newMainBranch {
					return errors.New(fmt.Sprintf("'%s' is already set as the main branch", branchName))
				}
				return nil
			},
		})
		if newPerennialBranch == "" {
			break
		}
		newPerennialBranches = append(newPerennialBranches, newPerennialBranch)
	}
	git.SetPerennialBranches(newPerennialBranches)
}

// Helpers

var configurationHeaderShown bool

func getMainBranchPrompt() (result string) {
	result += "Please specify the main development branch by name or number"
	currentMainBranch := git.GetMainBranch()
	if currentMainBranch != "" {
		coloredBranchName := color.New(color.Bold).Add(color.FgCyan).Sprintf(currentMainBranch)
		result += fmt.Sprintf(" (current value: %s)", coloredBranchName)
	}
	result += ":"
	return
}

func getPerennialBranchesPrompt() (result string) {
	result += "Please specify a perennial branch by name or number. Leave it blank to finish"
	currentPerennialBranches := git.GetPerennialBranches()
	if len(currentPerennialBranches) > 0 {
		coloredBranchNames := color.New(color.Bold).Add(color.FgCyan).Sprintf(strings.Join(currentPerennialBranches, ", "))
		result += fmt.Sprintf(" (current value: %s)", coloredBranchNames)
	}
	result += ":"
	return
}

func printConfigurationHeader() {
	if !configurationHeaderShown {
		configurationHeaderShown = true
		fmt.Println("Git Town needs to be configured\n")
		printNumberedBranches(git.GetLocalBranches())
		fmt.Println()
	}
}

// #!/usr/bin/env bash
//
//
// # Prints the header for the prompt when asking for configuration
// function echo_configuration_header {
//   echo "Git Town needs to be configured"
//   echo
//   echo_numbered_branches_alpha_order
//   echo
// }
//
//
// # Makes sure Git Town is configured
// function ensure_knows_configuration {
//   local header_shown=false
//   local numerical_regex='^[0-9]+$'
//   local user_input
//
//   if [ "$header_shown" = false ]; then
//     echo_configuration_header
//     header_shown=true
//   fi
//
//   local main_branch_input
//   local main_branch_current_value='none'
//   if [ "$(is_main_branch_configured)" = true ]; then
//     main_branch_current_value=$(echo_inline_cyan_bold "$MAIN_BRANCH_NAME")
//   fi
//
//   while [ -z "$main_branch_input" ]; do
//     echo -n "Please specify the main development branch by name or number (current value: $main_branch_current_value): "
//
//     read user_input
//     if [[ $user_input =~ $numerical_regex ]] ; then
//       main_branch_input="$(get_numbered_branch_alpha_order "$user_input")"
//       if [ -z "$main_branch_input" ]; then
//         echo_error_header
//         echo_error "Invalid branch number"
//         echo
//       fi
//     elif [ -z "$user_input" ]; then
//       if [ "$(is_main_branch_configured)" = true ]; then
//         main_branch_input=$MAIN_BRANCH_NAME
//       else
//         echo_error_header
//         echo_error "A main development branch is required to enable the features provided by Git Town"
//         echo
//       fi
//     else
//       if [ "$(has_branch "$user_input")" == true ]; then
//         main_branch_input=$user_input
//       else
//         echo_error_header
//         echo_error "Branch '$user_input' doesn't exist"
//         echo
//       fi
//     fi
//   done
//
//   store_configuration main-branch-name "$main_branch_input"
//
//
//   local perennial_branches_input=''
//   local perennial_branches_current_values='None'
//   if [ "$(are_perennial_branches_configured)" = true ]; then
//     perennial_branches_current_values=$(echo_inline_cyan_bold "$PERENNIAL_BRANCH_NAMES")
//   fi
//
//   while true; do
//     echo -n "Please specify a perennial branch by name or number. Leave it blank to finish (current value: $perennial_branches_current_values): "
//
//     read user_input
//     local branch
//     if [[ $user_input =~ $numerical_regex ]] ; then
//       branch="$(get_numbered_branch_alpha_order "$user_input")"
//       if [ -z "$branch" ]; then
//         echo_error_header
//         echo_error "Invalid branch number"
//         echo
//       fi
//     elif [ -z "$user_input" ]; then
//       break
//     else
//       if [ "$(has_branch "$user_input")" == true ]; then
//         if [ "$user_input" == "$MAIN_BRANCH_NAME" ]; then
//           echo_error_header
//           echo_error "'$user_input' is already set as the main branch"
//           echo
//         else
//           branch=$user_input
//         fi
//       else
//         echo_error_header
//         echo_error "Branch '$user_input' doesn't exist"
//         echo
//       fi
//     fi
//
//     if [ -n "$branch" ]; then
//       perennial_branches_input="$(insert_string "$perennial_branches_input" ' ' "$branch")"
//     fi
//   done
//
//   store_configuration perennial-branch-names "$perennial_branches_input"
// }
