package dialog

import (
	"fmt"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/git-town/git-town/v9/src/messages"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

const (
	// ResponseAbort stands for the user choosing to abort the unfinished run state.
	ResponseAbort = "abort"
	// ResponseContinue stands for the user choosing to continue the unfinished run state.
	ResponseContinue = "continue"
	// ResponseDiscard stands for the user choosing to discard the unfinished run state.
	ResponseDiscard = "discard"
	// ResponseQuit stands for the user choosing to quit the program.
	ResponseQuit = "quit"
	// ResponseSkip stands for the user choosing to continue the unfinished run state by skipping the current branch.
	ResponseSkip = "skip"
)

// AskHowToHandleUnfinishedRunState prompts the user for how to handle the unfinished run state.
func AskHowToHandleUnfinishedRunState(command, endBranch string, endTime time.Time, canSkip bool) (string, error) {
	formattedOptions := map[string]string{
		ResponseAbort:    fmt.Sprintf("Abort the `%s` command", command),
		ResponseContinue: fmt.Sprintf("Restart the `%s` command after having resolved conflicts", command),
		ResponseDiscard:  "Discard the unfinished state and run the new command",
		ResponseQuit:     "Quit without running anything",
		ResponseSkip:     fmt.Sprintf("Restart the `%s` command by skipping the current branch", command),
	}
	options := []string{
		formattedOptions[ResponseQuit],
		formattedOptions[ResponseContinue],
	}
	if canSkip {
		options = append(options, formattedOptions[ResponseSkip])
	}
	options = append(options, formattedOptions[ResponseAbort], formattedOptions[ResponseDiscard])
	prompt := &survey.Select{
		Message: fmt.Sprintf("You have an unfinished `%s` command that ended on the `%s` branch %s. Please choose how to proceed", command, endBranch, humanize.Time(endTime)),
		Options: options,
		Default: formattedOptions[ResponseQuit],
	}
	result := ""
	err := survey.AskOne(prompt, &result, nil)
	if err != nil {
		return "", fmt.Errorf(messages.DialogCannotReadAnswer, err)
	}
	for responseType, formattedResponseType := range formattedOptions {
		if formattedResponseType == result {
			return responseType, nil
		}
	}
	return "", fmt.Errorf(messages.DialogUnexpectedResponse, result)
}
