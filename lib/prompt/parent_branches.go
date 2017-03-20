package prompt

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/Originate/git-town/lib/config"
	"github.com/Originate/git-town/lib/git"
	"github.com/Originate/git-town/lib/util"
	"github.com/fatih/color"
)

func EnsureKnowsParentBranches(branchNames []string) {
	for _, branchName := range branchNames {
		if config.IsMainBranch(branchName) || config.IsPerennialBranch(branchName) || config.HasCompiledAncestorBranches(branchName) {
			continue
		}
		askForBranchAncestry(branchName)
		ancestors := config.CompileAncestorBranches(branchName)
		config.SetAncestorBranches(branchName, ancestors)

		if parentBranchHeaderShown {
			fmt.Println()
		}
	}
}

// Helpers

var parentBranchHeaderShown = false
var parentBranchHeaderTemplate = `
Feature branches can be branched directly off
%s or from other feature branches.

The former allows to develop and ship features completely independent of each other.
The latter allows to build on top of currently unshipped features.

`
var parentBranchPromptTemplate = "Please specify the parent branch of %s by name or number (default: %s): "

func askForBranchAncestry(branchName string) {
	current := branchName
	for {
		parent := config.GetParentBranch(current)
		if parent == "" {
			printParentBranchHeader()
			parent = askForParentBranch(current)
			config.SetParentBranch(current, parent)
		}
		if parent == config.GetMainBranch() || config.IsPerennialBranch(parent) {
			break
		}
		current = parent
	}
}

func askForParentBranch(branchName string) string {
	for {
		printParentBranchPrompt(branchName)
		parent := parseParentBranch(util.GetUserInput())
		if parent == "" {
			continue
		} else if branchName == parent {
			util.PrintError(fmt.Sprintf("'%s' cannot be the parent of itself", parent))
		} else if config.IsAncestorBranch(parent, branchName) {
			util.PrintError(fmt.Sprintf("Nested branch loop detected: '%s' is an ancestor of '%s'", branchName, parent))
		} else {
			return parent
		}
	}
}

func parseParentBranch(userInput string) string {
	mainBranch := config.GetMainBranch()
	numericRegex, err := regexp.Compile("^[0-9]+$")
	if err != nil {
		log.Fatal("Error compiling numeric regular expression: ", err)
	}

	if numericRegex.MatchString(userInput) {
		return parseParentBranchNumber(userInput)
	} else if userInput == "" {
		return mainBranch
	} else if git.HasBranch(userInput) {
		return userInput
	} else {
		util.PrintError(fmt.Sprintf("Branch '%s' doesn't exist", userInput))
	}

	return ""
}

func parseParentBranchNumber(userInput string) string {
	numberedBranches := git.GetLocalBranchesWithMainBranchFirst()
	index, err := strconv.Atoi(userInput)
	if err != nil {
		log.Fatal("Error parsing string to integer: ", err)
	}
	if index >= 1 && index <= len(numberedBranches) {
		return numberedBranches[index-1]
	} else {
		util.PrintError("Invalid branch number")
		return ""
	}
}

func printNumberedBranches() {
	boldFmt := color.New(color.Bold)
	branches := git.GetLocalBranchesWithMainBranchFirst()
	for index, branchName := range branches {
		fmt.Printf("  %s: %s\n", boldFmt.Sprintf("%d", index+1), branchName)
	}
}

func printParentBranchHeader() {
	if !parentBranchHeaderShown {
		parentBranchHeaderShown = true
		fmt.Printf(parentBranchHeaderTemplate, config.GetMainBranch())
		printNumberedBranches()
		fmt.Println()
	}
}

func printParentBranchPrompt(branchName string) {
	coloredBranchName := color.New(color.Bold).Add(color.FgCyan).Sprintf(branchName)
	fmt.Printf(parentBranchPromptTemplate, coloredBranchName, config.GetMainBranch())
}
