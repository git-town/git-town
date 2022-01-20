package dialog

import (
	"runtime"

	surveyCore "gopkg.in/AlecAivazis/survey.v1/core"
)

// init configures the prompts to work on Windows.
func init() {
	if runtime.GOOS == "windows" {
		surveyCore.SelectFocusIcon = ">"
		surveyCore.MarkedOptionIcon = "[x]"
		surveyCore.UnmarkedOptionIcon = "[ ]"
	}
}
