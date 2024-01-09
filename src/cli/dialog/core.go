// Package dialog allows the user to enter configuration data via CLI dialogs and prompts.
package dialog

import (
	"runtime"

	tea "github.com/charmbracelet/bubbletea"
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

// Typical colors used in BubbleTea dialogs.
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

// bubbleList contains common elements of BubbleTea list implementations.
type bubbleList struct {
	colors  dialogColors // colors to use for help text
	cursor  int          // index of the currently selected row
	entries []string     // the entries to select from
}

func (self *bubbleList) moveCursorDown() {
	if self.cursor < len(self.entries)-1 {
		self.cursor++
	} else {
		self.cursor = 0
	}
}

func (self *bubbleList) moveCursorUp() {
	if self.cursor > 0 {
		self.cursor--
	} else {
		self.cursor = len(self.entries) - 1
	}
}

func (self bubbleList) selectedEntry() string {
	return self.entries[self.cursor]
}

func (self *bubbleList) handleKey(key tea.KeyMsg) bool {
	switch key.Type {
	case tea.KeyUp, tea.KeyShiftTab:
		self.moveCursorUp()
		return true
	case tea.KeyDown, tea.KeyTab:
		self.moveCursorDown()
		return true
	}
	switch key.String() {
	case "k", "A", "Z":
		self.moveCursorUp()
		return true
	case "j", "B":
		self.moveCursorDown()
		return true
	}
	return false
}
