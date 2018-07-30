package prompt

import (
	"fmt"
	"log"
	"time"

	"github.com/Originate/exit"
	humanize "github.com/dustin/go-humanize"
	survey "gopkg.in/AlecAivazis/survey.v1"
)

var (
	// ResponseTypeAbort stands for the user choosing to abort the unfinished run state
	ResponseTypeAbort = "abort"
	// ResponseTypeContinue stands for the user choosing to continue the unfinished run state
	ResponseTypeContinue = "continue"
	// ResponseTypeDiscard stands for the user choosing to discard the unfinished run state
	ResponseTypeDiscard = "discard"
	// ResponseTypeQuit stands for the user choosing to quit the program
	ResponseTypeQuit = "quit"
	// ResponseTypeSkip stands for the user choosing to continue the unfinished run state by skipping the current branch
	ResponseTypeSkip = "skip"
)

// AskHowToHandleUnfinishedRunState prompts the user for how to handle the unfinished run state
func AskHowToHandleUnfinishedRunState(command, endBranch string, endTime time.Time, canSkip bool) string {
	formattedOptions := map[string]string{
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
	exit.If(err)
	for responseType, formattedResponseType := range formattedOptions {
		if formattedResponseType == result {
			return responseType
		}
	}
	log.Fatalf("Unexpected response: %s", result)
	return ""
}
