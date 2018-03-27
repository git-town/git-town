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
	ResponseTypeAbort    = "abort"
	ResponseTypeContinue = "continue"
	ResponseTypeDiscard  = "discard"
	ResponseTypeQuit     = "quit"
	ResponseTypeSkip     = "skip"
)

func AskHowToHandleUnfinishedRunState(command, endBranch string, endTime time.Time, canSkip bool) string {
	formattedOptions := map[string]string{
		ResponseTypeAbort:    fmt.Sprintf("Abort the `%s` command", command),
		ResponseTypeContinue: fmt.Sprintf("Continue the `%s` command", command),
		ResponseTypeDiscard:  "Discard the unfinished state and run the new command",
		ResponseTypeQuit:     "Quit without running anything.",
		ResponseTypeSkip:     fmt.Sprintf("Skip the `%s` command", command),
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
		Message: fmt.Sprintf("You have an unfinished `%s` command that ended on the `%s` branch %s.", command, endBranch, humanize.Time(endTime)),
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
