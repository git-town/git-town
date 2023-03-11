package dialog

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// EnterPerennialBranches lets the user enter the perennial branches.
func EnterPerennialBranches(localBranchesWithoutMain []string, perennialBranches []string) ([]string, error) {
	result := []string{}
	if len(localBranchesWithoutMain) == 0 {
		return result, nil
	}
	prompt := &survey.MultiSelect{
		Message: perennialBranchesPrompt(perennialBranches),
		Options: localBranchesWithoutMain,
		Default: perennialBranches,
	}
	err := survey.AskOne(prompt, &result, nil)
	if err != nil {
		return result, fmt.Errorf("cannot read branches from CLI: %w", err)
	}
	return result, nil
}

func perennialBranchesPrompt(perennialBranches []string) string {
	result := "Please specify perennial branches:"
	if len(perennialBranches) > 0 {
		coloredBranches := color.New(color.Bold).Add(color.FgCyan).Sprintf(strings.Join(perennialBranches, ", "))
		result += fmt.Sprintf(" (current value: %s)", coloredBranches)
	}
	return result
}
