package dialog

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// EnterPerennialBranches lets the user enter the perennial branches.
func EnterPerennialBranches(localBranchesWithoutMain []string, perennialBranches []string) ([]string, error) {
	if len(localBranchesWithoutMain) == 0 {
		return []string{}, nil
	}
	newPerennialBranches, err := askForBranches(askForBranchesOptions{
		branches:        localBranchesWithoutMain,
		prompt:          perennialBranchesPrompt(perennialBranches),
		defaultBranches: perennialBranches,
	})
	if err != nil {
		return []string{}, err
	}
	return newPerennialBranches, nil
}

// Helpers

func askForBranches(opts askForBranchesOptions) ([]string, error) {
	result := []string{}
	prompt := &survey.MultiSelect{
		Message: opts.prompt,
		Options: opts.branches,
		Default: opts.defaultBranches,
	}
	err := survey.AskOne(prompt, &result, nil)
	if err != nil {
		return result, fmt.Errorf("cannot read branches from CLI: %w", err)
	}
	return result, err
}

type askForBranchesOptions struct {
	branches        []string
	defaultBranches []string
	prompt          string
}

func perennialBranchesPrompt(perennialBranches []string) string {
	result := "Please specify perennial branches:"
	if len(perennialBranches) > 0 {
		coloredBranches := color.New(color.Bold).Add(color.FgCyan).Sprintf(strings.Join(perennialBranches, ", "))
		result += fmt.Sprintf(" (current value: %s)", coloredBranches)
	}
	return result
}
