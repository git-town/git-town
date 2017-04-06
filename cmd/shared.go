package cmd

import "errors"

// These variables represent command-line flags
var (
	AbortFlag    bool
	AllFlag      bool
	ContinueFlag bool
	SkipFlag     bool
	UndoFlag     bool
)

func validateMaxArgs(args []string, max int) error {
	if len(args) > max {
		return errors.New("Too many arguments")
	}
	return nil
}
