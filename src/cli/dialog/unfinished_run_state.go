package dialog

import (
	"fmt"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/messages"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

type Response struct {
	name string
}

func (self Response) String() string { return self.name }

var (
	// ResponseContinue stands for the user choosing to continue the unfinished run state.
	ResponseContinue = Response{"continue"} //nolint:gochecknoglobals
	// ResponseDiscard stands for the user choosing to discard the unfinished run state.
	ResponseDiscard = Response{"discard"} //nolint:gochecknoglobals
	// ResponseQuit stands for the user choosing to quit the program.
	ResponseQuit = Response{"quit"} //nolint:gochecknoglobals
	// ResponseSkip stands for the user choosing to continue the unfinished run state by skipping the current branch.
	ResponseSkip = Response{"skip"} //nolint:gochecknoglobals
	// ResponseUndo stands for the user choosing to undo the unfinished run state.
	ResponseUndo = Response{"undo"} //nolint:gochecknoglobals
)

// AskHowToHandleUnfinishedRunState prompts the user for how to handle the unfinished run state.
func AskHowToHandleUnfinishedRunState(command string, endBranch domain.LocalBranchName, endTime time.Time, canSkip bool) (Response, error) {
	formattedOptions := map[Response]string{
		ResponseContinue: fmt.Sprintf("Restart the `%s` command after having resolved conflicts", command),
		ResponseDiscard:  "Discard the unfinished state and run the new command",
		ResponseQuit:     "Quit without running anything",
		ResponseSkip:     fmt.Sprintf("Restart the `%s` command by skipping the current branch", command),
		ResponseUndo:     fmt.Sprintf("Undo the `%s` command", command), // TODO: extract these strings to en.go.
	}
	options := []string{
		formattedOptions[ResponseQuit],
		formattedOptions[ResponseContinue],
	}
	if canSkip {
		options = append(options, formattedOptions[ResponseSkip])
	}
	options = append(options, formattedOptions[ResponseUndo], formattedOptions[ResponseDiscard])
	prompt := &survey.Select{
		Message: fmt.Sprintf("You have an unfinished `%s` command that ended on the `%s` branch %s. Please choose how to proceed", command, endBranch, humanize.Time(endTime)),
		Options: options,
		Default: formattedOptions[ResponseQuit],
	}
	result := ""
	err := survey.AskOne(prompt, &result, nil)
	if err != nil {
		return ResponseUndo, fmt.Errorf(messages.DialogCannotReadAnswer, err)
	}
	for responseType, formattedResponseType := range formattedOptions {
		if formattedResponseType == result {
			return responseType, nil
		}
	}
	return ResponseUndo, fmt.Errorf(messages.DialogUnexpectedResponse, result)
}
