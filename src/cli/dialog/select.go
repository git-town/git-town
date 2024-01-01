package dialog

import (
	"fmt"

	"github.com/git-town/git-town/v11/src/messages"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

// Select displays a visual dialog that allows the user to select one of the given options.
func Select(opts SelectArgs) (string, error) {
	result := ""
	prompt := &survey.Select{
		Message: opts.Message,
		Options: opts.Options,
		Default: opts.Default,
	}
	err := survey.AskOne(prompt, &result, nil)
	if err != nil {
		return result, fmt.Errorf(messages.DialogCannotReadBranch, err)
	}
	return result, nil
}

type SelectArgs struct {
	Options []string
	Default string
	Message string
}
