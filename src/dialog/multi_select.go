package dialog

import survey "gopkg.in/AlecAivazis/survey.v1"

func MultiSelect(args MultiSelectArgs) ([]string, error) {
	result := []string{}
	if len(args.Options) == 0 {
		return result, nil
	}
	prompt := &survey.MultiSelect{
		Message: args.Message,
		Options: args.Options,
		Default: args.Default,
	}
	err := survey.AskOne(prompt, &result, nil)
	return result, err
}

type MultiSelectArgs struct {
	Options []string
	Default []string
	Message string
}
