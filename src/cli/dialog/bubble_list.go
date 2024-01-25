package dialog

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/gohacks"
	"github.com/muesli/termenv"
)

type dialogStatus int

const (
	dialogStatusActive  dialogStatus = iota // the user is currently entering data into the dialog
	dialogStatusDone                        // the user has made a selection
	dialogStatusAborted                     // the user has aborted the dialog
)

// BubbleList contains common elements of BubbleTea list implementations.
type BubbleList[S fmt.Stringer] struct {
	Status       dialogStatus
	Colors       dialogColors  // colors to use for help text
	Cursor       int           // index of the currently selected row
	Dim          termenv.Style // style for dim output
	Entries      []S           // the entries to select from
	EntryNumber  string        // the manually entered entry number
	MaxDigits    int           // how many digits make up an entry number
	NumberFormat string        // template for formatting the entry number
}

func newBubbleList[S fmt.Stringer](entries []S, cursor int) BubbleList[S] {
	numberLen := gohacks.NumberLength(len(entries))
	return BubbleList[S]{
		Status:       dialogStatusActive,
		Colors:       createColors(),
		Cursor:       cursor,
		Dim:          termenv.String().Faint(),
		Entries:      entries,
		EntryNumber:  "",
		MaxDigits:    numberLen,
		NumberFormat: fmt.Sprintf("%%0%dd ", numberLen),
	}
}

// Aborted indicates whether the user has aborted this dialog.
func (self *BubbleList[S]) aborted() bool {
	return self.Status == dialogStatusAborted
}

// entryNumberStr provides a colorized string to print the given entry number.
func (self *BubbleList[S]) entryNumberStr(number int) string {
	return self.Dim.Styled(fmt.Sprintf(self.NumberFormat, number))
}

// handleKey handles keypresses that are common for all bubbleLists.
func (self *BubbleList[S]) handleKey(key tea.KeyMsg) (bool, tea.Cmd) {
	switch key.Type { //nolint:exhaustive
	case tea.KeyUp, tea.KeyShiftTab:
		self.moveCursorUp()
		return true, nil
	case tea.KeyDown, tea.KeyTab:
		self.moveCursorDown()
		return true, nil
	case tea.KeyCtrlC, tea.KeyEsc:
		self.Status = dialogStatusAborted
		return true, tea.Quit
	}
	switch keyStr := key.String(); keyStr {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		self.EntryNumber += keyStr
		if len(self.EntryNumber) > self.MaxDigits {
			self.EntryNumber = self.EntryNumber[1:]
		}
		number64, _ := strconv.ParseInt(self.EntryNumber, 10, 0)
		number := int(number64)
		if number < len(self.Entries) {
			self.Cursor = number
		}
	case "k", "A", "Z":
		self.moveCursorUp()
		return true, nil
	case "j", "B":
		self.moveCursorDown()
		return true, nil
	case "q":
		self.Status = dialogStatusAborted
		return true, tea.Quit
	}
	return false, nil
}

func (self *BubbleList[S]) moveCursorDown() {
	if self.Cursor < len(self.Entries)-1 {
		self.Cursor++
	} else {
		self.Cursor = 0
	}
}

func (self *BubbleList[S]) moveCursorUp() {
	if self.Cursor > 0 {
		self.Cursor--
	} else {
		self.Cursor = len(self.Entries) - 1
	}
}

func (self BubbleList[S]) selectedEntry() S { //nolint:ireturn
	return self.Entries[self.Cursor]
}
