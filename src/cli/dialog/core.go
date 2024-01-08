// Package dialog allows the user to enter configuration data via CLI dialogs and prompts.
package dialog

import (
	"runtime"

	"github.com/muesli/termenv"
	surveyCore "gopkg.in/AlecAivazis/survey.v1/core"
)

// Initialize configures the prompts to work on Windows.
func Initialize() {
	if runtime.GOOS == "windows" {
		surveyCore.SelectFocusIcon = ">"
		surveyCore.MarkedOptionIcon = "[x]"
		surveyCore.UnmarkedOptionIcon = "[ ]"
	}
}

// Typical colors used in BubbleTea dialogs
type dialogColors struct {
	help      termenv.Style // color of help text
	helpKey   termenv.Style // color of key names in help text
	initial   termenv.Style // color for the row containing the currently checked out branch
	selection termenv.Style // color for the currently selected entry
}

func createColors() dialogColors {
	return dialogColors{
		help:      termenv.String().Faint(),
		helpKey:   termenv.String().Faint().Bold(),
		initial:   termenv.String().Foreground(termenv.ANSIGreen),
		selection: termenv.String().Foreground(termenv.ANSICyan),
	}
}
