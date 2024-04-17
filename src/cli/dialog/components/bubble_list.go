package components

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v14/src/cli/colors"
	"github.com/git-town/git-town/v14/src/gohacks"
	"github.com/muesli/termenv"
)

type status int

const (
	StatusActive  status = iota // the user is currently entering data into the dialog
	StatusDone                  // the user has made a selection
	StatusAborted               // the user has aborted the dialog
)

type BubbleListEntry[S fmt.Stringer] struct {
	// TODO: Checked bool
	Data    S
	Enabled bool
	Text    string
}

type BubbleListEntries[S fmt.Stringer] []BubbleListEntry[S]

func (self BubbleListEntries[S]) allDisabled() bool {
	for _, entry := range self {
		if entry.Enabled {
			return false
		}
	}
	return true
}

// BubbleList contains common elements of BubbleTea list implementations.
type BubbleList[S fmt.Stringer] struct {
	Colors       dialogColors         // colors to use for help text
	Cursor       int                  // index of the currently selected row
	Dim          termenv.Style        // style for dim output
	Entries      BubbleListEntries[S] // the entries to select from
	EntryNumber  string               // the manually entered entry number
	MaxDigits    int                  // how many digits make up an entry number
	NumberFormat string               // template for formatting the entry number
	Status       status
}

func NewBubbleList[S fmt.Stringer](entries []BubbleListEntry[S], cursor int) BubbleList[S] {
	numberLen := gohacks.NumberLength(len(entries))
	return BubbleList[S]{
		Status:       StatusActive,
		Colors:       createColors(),
		Cursor:       cursor,
		Dim:          colors.Faint(),
		Entries:      entries,
		EntryNumber:  "",
		MaxDigits:    numberLen,
		NumberFormat: fmt.Sprintf("%%0%dd ", numberLen),
	}
}

// NewEnabledBubbleListEntries creates enabled BubbleListEntries for the given data types.
func NewEnabledBubbleListEntries[S fmt.Stringer](records []S) []BubbleListEntry[S] {
	result := make([]BubbleListEntry[S], len(records))
	for r, record := range records {
		result[r] = BubbleListEntry[S]{
			Data:    record,
			Enabled: true,
			Text:    record.String(),
		}
	}
	return result
}

// Aborted indicates whether the user has Aborted this components.
func (self *BubbleList[S]) Aborted() bool {
	return self.Status == StatusAborted
}

// EntryNumberStr provides a colorized string to print the given entry number.
func (self *BubbleList[S]) EntryNumberStr(number int) string {
	return self.Dim.Styled(fmt.Sprintf(self.NumberFormat, number))
}

// HandleKey handles keypresses that are common for all bubbleLists.
func (self *BubbleList[S]) HandleKey(key tea.KeyMsg) (bool, tea.Cmd) {
	switch key.Type { //nolint:exhaustive
	case tea.KeyUp, tea.KeyShiftTab:
		self.moveCursorUp()
		return true, nil
	case tea.KeyDown, tea.KeyTab:
		self.moveCursorDown()
		return true, nil
	case tea.KeyLeft:
		self.movePageUp()
		return true, nil
	case tea.KeyRight:
		self.movePageDown()
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
		self.moveCursorUp()
		return true, nil
	case "j":
		self.moveCursorDown()
		return true, nil
	case "u":
		self.movePageUp()
		return true, nil
	case "d":
		self.movePageDown()
		return true, nil
	case "q":
		self.Status = StatusAborted
		return true, tea.Quit
	}
	return false, nil
}

func (self BubbleList[S]) SelectedEntry() BubbleListEntry[S] { //nolint:ireturn
	return self.Entries[self.Cursor]
}

func (self BubbleList[S]) SelectedData() S { //nolint:ireturn
	return self.SelectedEntry().Data
}

func (self *BubbleList[S]) moveCursorDown() {
	if self.Entries.allDisabled() {
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

func (self *BubbleList[S]) moveCursorUp() {
	if self.Entries.allDisabled() {
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

func (self *BubbleList[S]) movePageDown() {
	if self.Entries.allDisabled() {
		return
	}
	self.Cursor += 10
	if self.Cursor >= len(self.Entries) {
		self.Cursor = len(self.Entries) - 1
	}
	if self.SelectedEntry().Enabled {
		return
	}
	// keep moving the cursor down until we hit either an enabled entry or the end of the list
	for self.Cursor < len(self.Entries)-1 {
		self.Cursor++
		if self.SelectedEntry().Enabled {
			return
		}
	}
	// here we have reached the end of the list and haven't found an enabled entry --> go backwards until we find one
	self.moveCursorUp()
}

func (self *BubbleList[S]) movePageUp() {
	if self.Entries.allDisabled() {
		return
	}
	self.Cursor -= 10
	if self.Cursor < 0 {
		self.Cursor = 0
	}
	if self.SelectedEntry().Enabled {
		return
	}
	// keep moving the cursor up until we hit either an enabled entry or the beginning of the list
	for self.Cursor > 0 {
		self.Cursor--
		if self.SelectedEntry().Enabled {
			return
		}
	}
	// here we have reached the end of the list and haven't found an enabled entry --> go backwards until we find one
	self.moveCursorDown()
}
