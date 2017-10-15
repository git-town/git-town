package prompt

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/Originate/exit"
	"github.com/Originate/git-town/src/cfmt"
	"github.com/Originate/git-town/src/git"
	"github.com/Originate/git-town/src/util"
	"github.com/fatih/color"
)

type branchPromptConfig struct {
	branchNames       []string
	defaultBranchName string
	prompt            string
	validate          func(branchName string) error
}

func askForBranch(config branchPromptConfig) string {
	for {
		cfmt.Print(config.prompt)
		branchName, err := parseBranch(config, util.GetUserInput())
		if err == nil {
			err = config.validate(branchName)
			if err == nil {
				return branchName
			}
		}
		util.PrintError(err.Error())
	}
}

func parseBranch(config branchPromptConfig, userInput string) (string, error) {
	numericRegex, err := regexp.Compile("^[0-9]+$")
	exit.IfWrap(err, "Error compiling numeric regular expression")

	if numericRegex.MatchString(userInput) {
		return parseBranchNumber(config.branchNames, userInput)
	}
	if userInput == "" {
		return config.defaultBranchName, nil
	}
	if git.HasBranch(userInput) {
		return userInput, nil
	}

	return "", fmt.Errorf("Branch '%s' doesn't exist", userInput)
}

func parseBranchNumber(branchNames []string, userInput string) (string, error) {
	index, err := strconv.Atoi(userInput)
	exit.IfWrap(err, "Error parsing string to integer")
	if index >= 1 && index <= len(branchNames) {
		return branchNames[index-1], nil
	}

	return "", errors.New("Invalid branch number")
}

func printNumberedBranches(branchNames []string) {
	boldFmt := color.New(color.Bold)
	for index, branchName := range branchNames {
		cfmt.Printf("  %s: %s\n", boldFmt.Sprintf("%d", index+1), branchName)
	}
}
