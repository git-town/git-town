package dialog

import survey "gopkg.in/AlecAivazis/survey.v1"

// MultiSelect displays a visual dialog that allows the user to select multiple entries amongst the given options.
func MultiSelect(args MultiSelectArgs) ([]string, error) {
	result := []string{}
	if len(args.Options) == 0 {
		return result, nil
	}
	prompt := &survey.MultiSelect{
		Message: args.Message,
		Options: args.Options,
		Default: args.Defaults,
	}
	err := survey.AskOne(prompt, &result, nil)
	return result, err
}

type MultiSelectArgs struct {
	Options  []string
	Defaults []string
	Message  string
}
