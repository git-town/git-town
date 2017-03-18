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
	branchNameFmt := color.New(color.Bold).Add(color.FgCyan)
	branches := git.GetLocalBranchesWithMainBranchFirst()
	mainBranch := config.GetMainBranch()
	message := fmt.Sprintf("Please specify the parent branch of %s by name or number (default: %s): ", branchNameFmt.Sprintf(branchName), mainBranch)
	numericRegex, err := regexp.Compile("^[0-9]+$")
	if err != nil {
		log.Fatal("Error compiling numeric regular expression: ", err)
	}

	parent := ""
	for {
		fmt.Printf(message)
		userInput := util.GetUserInput()
		if numericRegex.MatchString(userInput) {
			index, err := strconv.Atoi(userInput)
			if err != nil {
				log.Fatal("Error parsing string to integer: ", err)
			}
			if index >= 1 && index <= len(branches) {
				parent = branches[index-1]
			} else {
				util.PrintError("Invalid branch number")
			}
		} else if userInput == "" {
			parent = mainBranch
		} else if git.HasBranch(userInput) {
			parent = userInput
		} else {
			util.PrintError(fmt.Sprintf("Branch '%s' doesn't exist", userInput))
		}

		if parent == "" {
			continue
		} else if branchName == parent {
			util.PrintError(fmt.Sprintf("'%s' cannot be the parent of itself", parent))
			parent = ""
		} else if config.HasAncestorBranch(parent, branchName) {
			util.PrintError(fmt.Sprintf("Nested branch loop detected: '%s' is an ancestor of '%s'", branchName, parent))
			parent = ""
		} else {
			return parent
		}
	}
}

var parentBranchHeaderShown = false

func printParentBranchHeader() {
	if !parentBranchHeaderShown {
		parentBranchHeaderShown = true
		fmt.Println()
		fmt.Println("Feature branches can be branched directly off ")
		fmt.Println("%s or from other feature branches.", config.GetMainBranch())
		fmt.Println()
		fmt.Println("The former allows to develop and ship features completely independent of each other.")
		fmt.Println("The latter allows to build on top of currently unshipped features.")
		fmt.Println()
		printNumberedBranches()
		fmt.Println()
	}
}

func printNumberedBranches() {
	boldFmt := color.New(color.Bold)
	branches := git.GetLocalBranchesWithMainBranchFirst()
	for index, branchName := range branches {
		fmt.Printf("  %s: %s\n", boldFmt.Sprintf("%d", index+1), branchName)
	}
}
