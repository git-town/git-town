package dialog

import (
	"slices"

	tea "github.com/charmbracelet/bubbletea"
)

// bubbleList contains common elements of BubbleTea list implementations.
type bubbleList struct {
	aborted bool         // whether the user has aborted this dialog
	colors  dialogColors // colors to use for help text
	cursor  int          // index of the currently selected row
	entries []string     // the entries to select from
}

func newBubbleList(entries []string, initial string) bubbleList {
	cursor := slices.Index(entries, initial)
	if cursor < 0 {
		cursor = 0
	}
	return bubbleList{
		aborted: false,
		entries: entries,
		colors:  createColors(),
		cursor:  cursor,
	}
}

func (self *bubbleList) handleKey(key tea.KeyMsg) (bool, tea.Cmd) {
	switch key.Type { //nolint:exhaustive
	case tea.KeyUp, tea.KeyShiftTab:
		self.moveCursorUp()
		return true, nil
	case tea.KeyDown, tea.KeyTab:
		self.moveCursorDown()
		return true, nil
	case tea.KeyCtrlC:
		self.aborted = true
		return true, tea.Quit
	}
	switch key.String() {
	case "k", "A", "Z":
		self.moveCursorUp()
		return true, nil
	case "j", "B":
		self.moveCursorDown()
		return true, nil
	case "q":
		self.aborted = true
		return true, tea.Quit
	}
	return false, nil
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
