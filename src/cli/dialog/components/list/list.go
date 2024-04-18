package list

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v14/src/cli/colors"
	"github.com/git-town/git-town/v14/src/gohacks"
)

// List contains elements and operations common to all BubbleTea-based list implementations.
type List[S fmt.Stringer] struct {
	Colors       colors.DialogColors // colors to use for help text
	Cursor       int                 // index of the currently selected row
	Entries      Entries[S]          // the entries to select from
	EntryNumber  string              // the manually entered entry number
	MaxDigits    int                 // how many digits make up an entry number
	NumberFormat string              // template for formatting the entry number
	Status       Status
}

func NewList[S fmt.Stringer](entries Entries[S], cursor int) List[S] {
	numberLen := gohacks.NumberLength(len(entries))
	return List[S]{
		Status:       StatusActive,
		Colors:       colors.NewDialogColors(),
		Cursor:       cursor,
		Entries:      entries,
		EntryNumber:  "",
		MaxDigits:    numberLen,
		NumberFormat: fmt.Sprintf("%%0%dd ", numberLen),
	}
}

// Aborted indicates whether the user has Aborted this components.
func (self *List[S]) Aborted() bool {
	return self.Status == StatusAborted
}

// EntryNumberStr provides a colorized string to print the given entry number.
func (self *List[S]) EntryNumberStr(number int) string {
	return self.Colors.EntryNumber.Styled(fmt.Sprintf(self.NumberFormat, number))
}

// HandleKey handles keypresses that are common for all bubbleLists.
func (self *List[S]) HandleKey(key tea.KeyMsg) (bool, tea.Cmd) {
	switch key.Type { //nolint:exhaustive
	case tea.KeyUp, tea.KeyShiftTab:
		self.MoveCursorUp()
		return true, nil
	case tea.KeyDown, tea.KeyTab:
		self.MoveCursorDown()
		return true, nil
	case tea.KeyLeft:
		self.MovePageUp()
		return true, nil
	case tea.KeyRight:
		self.MovePageDown()
		return true, nil
	case tea.KeyCtrlC, tea.KeyEsc:
		self.Status = StatusAborted
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
	case "k":
		self.MoveCursorUp()
		return true, nil
	case "j":
		self.MoveCursorDown()
		return true, nil
	case "u":
		self.MovePageUp()
		return true, nil
	case "d":
		self.MovePageDown()
		return true, nil
	case "q":
		self.Status = StatusAborted
		return true, tea.Quit
	}
	return false, nil
}

func (self *List[S]) MoveCursorDown() {
	if self.Entries.AllDisabled() {
		return
	}
	for {
		if self.Cursor < len(self.Entries)-1 {
			self.Cursor++
		} else {
			self.Cursor = 0
		}
		if self.SelectedEntry().Enabled {
			return
		}
	}
}

func (self *List[S]) MoveCursorUp() {
	if self.Entries.AllDisabled() {
		return
	}
	for {
		if self.Cursor > 0 {
			self.Cursor--
		} else {
			self.Cursor = len(self.Entries) - 1
		}
		if self.SelectedEntry().Enabled {
			return
		}
	}
}

func (self *List[S]) MovePageDown() {
	if self.Entries.AllDisabled() {
		return
	}
	self.Cursor += 10
	if self.Cursor >= len(self.Entries) {
		self.Cursor = len(self.Entries) - 1
	}
	// go down until we find a selected entry
	for self.Cursor < len(self.Entries)-1 {
		if self.SelectedEntry().Enabled {
			return
		}
		self.Cursor += 1
	}
	// go up until we find a selected entry
	for {
		if self.SelectedEntry().Enabled {
			return
		}
		self.Cursor -= 1
	}
}

func (self *List[S]) MovePageUp() {
	self.Cursor -= 10
	if self.Cursor < 0 {
		self.Cursor = 0
	}
}

func (self List[S]) SelectedData() S { //nolint:ireturn
	return self.SelectedEntry().Data
}

func (self List[S]) SelectedEntry() Entry[S] {
	return self.Entries[self.Cursor]
}
