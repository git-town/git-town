package dialog

import (
	"fmt"

	survey "gopkg.in/AlecAivazis/survey.v1"
)

type askForBranchOptions struct {
	branches      []string
	defaultBranch string
	prompt        string
}

type askForBranchesOptions struct {
	branches        []string
	defaultBranches []string
	prompt          string
}

func askForBranch(opts askForBranchOptions) (string, error) {
	result := ""
	prompt := &survey.Select{
		Message: opts.prompt,
		Options: opts.branches,
		Default: opts.defaultBranch,
	}
	err := survey.AskOne(prompt, &result, nil)
	if err != nil {
		return result, fmt.Errorf("cannot read branch from CLI: %w", err)
	}
	return result, nil
}

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
