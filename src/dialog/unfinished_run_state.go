package dialog

import (
	"fmt"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/git-town/git-town/v9/src/messages"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type ResponseType struct {
	name string
}

func (r ResponseType) String() string { return r.name }

var (
	// ResponseTypeAbort stands for the user choosing to abort the unfinished run state.
	ResponseTypeAbort = ResponseType{"abort"} //nolint:gochecknoglobals
	// ResponseTypeContinue stands for the user choosing to continue the unfinished run state.
	ResponseTypeContinue = ResponseType{"continue"} //nolint:gochecknoglobals
	// ResponseTypeDiscard stands for the user choosing to discard the unfinished run state.
	ResponseTypeDiscard = ResponseType{"discard"} //nolint:gochecknoglobals
	// ResponseTypeQuit stands for the user choosing to quit the program.
	ResponseTypeQuit = ResponseType{"quit"} //nolint:gochecknoglobals
	// ResponseTypeSkip stands for the user choosing to continue the unfinished run state by skipping the current branch.
	ResponseTypeSkip = ResponseType{"skip"} //nolint:gochecknoglobals
)

// AskHowToHandleUnfinishedRunState prompts the user for how to handle the unfinished run state.
func AskHowToHandleUnfinishedRunState(command, endBranch string, endTime time.Time, canSkip bool) (ResponseType, error) {
	formattedOptions := map[ResponseType]string{
		ResponseTypeAbort:    fmt.Sprintf("Abort the `%s` command", command),
		ResponseTypeContinue: fmt.Sprintf("Restart the `%s` command after having resolved conflicts", command),
		ResponseTypeDiscard:  "Discard the unfinished state and run the new command",
		ResponseTypeQuit:     "Quit without running anything",
		ResponseTypeSkip:     fmt.Sprintf("Restart the `%s` command by skipping the current branch", command),
	}
	options := []string{
		formattedOptions[ResponseTypeQuit],
		formattedOptions[ResponseTypeContinue],
	}
	if canSkip {
		options = append(options, formattedOptions[ResponseTypeSkip])
	}
	options = append(options, formattedOptions[ResponseTypeAbort], formattedOptions[ResponseTypeDiscard])
	prompt := &survey.Select{
		Message: fmt.Sprintf("You have an unfinished `%s` command that ended on the `%s` branch %s. Please choose how to proceed", command, endBranch, humanize.Time(endTime)),
		Options: options,
		Default: formattedOptions[ResponseTypeQuit],
	}
	result := ""
	err := survey.AskOne(prompt, &result, nil)
	if err != nil {
		return ResponseTypeAbort, fmt.Errorf(messages.DialogCannotReadAnswer, err)
	}
	for responseType, formattedResponseType := range formattedOptions {
		if formattedResponseType == result {
			return responseType, nil
		}
	}
	return ResponseTypeAbort, fmt.Errorf(messages.DialogUnexpectedResponse, result)
}
