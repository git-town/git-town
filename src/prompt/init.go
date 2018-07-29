package prompt

import (
	"runtime"

	surveyCore "gopkg.in/AlecAivazis/survey.v1/core"
)

// initPrompts configures the prompts to work on Windows
func initPrompts() {
	if runtime.GOOS == "windows" {
		surveyCore.SelectFocusIcon = ">"
		surveyCore.MarkedOptionIcon = "[x]"
		surveyCore.UnmarkedOptionIcon = "[ ]"
	}

}
