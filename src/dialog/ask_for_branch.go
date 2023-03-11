package dialog

import (
	"fmt"

	survey "gopkg.in/AlecAivazis/survey.v1"
)

func AskForBranch(opts AskForBranchOptions) (string, error) {
	result := ""
	prompt := &survey.Select{
		Message: opts.Prompt,
		Options: opts.Branches,
		Default: opts.DefaultBranch,
	}
	err := survey.AskOne(prompt, &result, nil)
	if err != nil {
		return result, fmt.Errorf("cannot read branch from CLI: %w", err)
	}
	return result, nil
}

type AskForBranchOptions struct {
	Branches      []string
	DefaultBranch string
	Prompt        string
}
