// Package dialog allows the user to enter configuration data via CLI dialogs and prompts.
package dialog

import (
	"runtime"

	surveyCore "gopkg.in/AlecAivazis/survey.v1/core"
)

// init configures the prompts to work on Windows.
func Initialize() {
	if runtime.GOOS == "windows" {
		surveyCore.SelectFocusIcon = ">"
		surveyCore.MarkedOptionIcon = "[x]"
		surveyCore.UnmarkedOptionIcon = "[ ]"
	}
}
