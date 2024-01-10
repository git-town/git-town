package dialog

import (
	"fmt"
	"slices"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/gohacks"
	"github.com/muesli/termenv"
)

// bubbleList contains common elements of BubbleTea list implementations.
type bubbleList struct {
	aborted      bool          // whether the user has aborted this dialog
	colors       dialogColors  // colors to use for help text
	cursor       int           // index of the currently selected row
	dim          termenv.Style // for dim output
	entries      []string      // the entries to select from
	entryNumber  string        // the currently entered branch number
	maxDigits    int           // the maximal number of digits in the branch number
	numberFormat string        // template for formatting the entry number
}

func newBubbleList(entries []string, initial string) bubbleList {
	cursor := slices.Index(entries, initial)
	if cursor < 0 {
		cursor = 0
	}
	numberLen := gohacks.NumberLength(len(entries))
	return bubbleList{
		aborted:      false,
		colors:       createColors(),
		cursor:       cursor,
		dim:          termenv.String().Faint(),
		entries:      entries,
		entryNumber:  "",
		maxDigits:    numberLen,
		numberFormat: fmt.Sprintf("%%0%dd ", numberLen),
	}
}

func (self *bubbleList) entryNumberStr(number int) string {
	return self.dim.Styled(fmt.Sprintf(self.numberFormat, number))
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
	switch keyStr := key.String(); keyStr {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		self.entryNumber += keyStr
		if len(self.entryNumber) > self.maxDigits {
			self.entryNumber = self.entryNumber[1:]
		}
		number64, _ := strconv.ParseInt(self.entryNumber, 10, 0)
		number := int(number64)
		if number < len(self.entries) {
			self.cursor = number
		}
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
