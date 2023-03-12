package dialog

import (
	"fmt"

	survey "gopkg.in/AlecAivazis/survey.v1"
)

// Select allows the user to select one of the given branches.
func Select(opts SelectArgs) (string, error) {
	result := ""
	prompt := &survey.Select{
		Message: opts.Message,
		Options: opts.Options,
		Default: opts.Default,
	}
	err := survey.AskOne(prompt, &result, nil)
	if err != nil {
		return result, fmt.Errorf("cannot read branch from CLI: %w", err)
	}
	return result, nil
}

type SelectArgs struct {
	Options []string
	Default string
	Message string
}
