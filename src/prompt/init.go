package prompt

import (
	"runtime"

	surveyCore "gopkg.in/AlecAivazis/survey.v1/core"
)

// Init configures the prompts to work on Windows
func Init() {
	if runtime.GOOS == "windows" {
		surveyCore.SelectFocusIcon = ">"
		surveyCore.MarkedOptionIcon = "[x]"
		surveyCore.UnmarkedOptionIcon = "[ ]"
	}
}
