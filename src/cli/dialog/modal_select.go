package dialog

import (
	"fmt"
	"strings"

	"atomicgo.dev/cursor"
	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/git-town/git-town/v11/src/messages"
)

// ModalSelect allows the user to select a value from the given entries.
// Entries can be arbitrarily formatted.
// The given initial value is preselected.
func ModalSelect(entries ModalSelectEntries, initialValue string) (*string, error) {
	initialPos := entries.IndexOfValue(initialValue)
	if initialPos == nil {
		return nil, fmt.Errorf(messages.DialogOptionNotFound, initialValue, entries)
	}
	input := modalSelect{
		entries:       entries,
		activeCursor:  "> ",
		activeColor:   color.New(color.FgCyan, color.Bold),
		activePos:     *initialPos,
		initialCursor: "* ",
		initialColor:  color.New(color.FgGreen),
		initialPos:    *initialPos,
		status:        modalSelectStatusNew,
	}
	return input.Display()
}

// modalSelect allows selecting a value from a list using VIM keybindings.
type modalSelect struct {
	activeColor   *color.Color       // color with which to print the currently selected line
	activeCursor  string             // text that gets prepended to the currently selected row
	activePos     int                // index of the currently selected row
	entries       ModalSelectEntries // the entries to display
	initialColor  *color.Color       // color with which to print the initially selected value
	initialCursor string             // cursor at the initial entry
	initialPos    int                // index of the initially selected value
	status        modalSelectStatus  // the current status of this ModalInput instance
}

// Display shows the dialog and lets the user select an entry.
// Returns the selected value or nil if the user aborted the dialog.
func (self *modalSelect) Display() (*string, error) {
	cursor.Hide()
	defer cursor.Show()
	err := keyboard.Open()
	if err != nil {
		return nil, err
	}
	defer keyboard.Close()
	self.print()
	for self.status == modalSelectStatusSelecting {
		err := self.handleInput()
		if err != nil {
			return nil, err
		}
		self.print()
	}
	if self.status == modalSelectStatusAborted {
		return nil, nil //nolint:nilnil
	}
	selectedValue := self.selectedValue()
	return &selectedValue, nil
}

func (self modalSelectStatus) String() string { return self.name }

// handleInput waits for keyboard input and updates the dialog state.
func (self *modalSelect) handleInput() error {
	char, key, err := keyboard.GetSingleKey()
	if err != nil {
		return err
	}
	switch {
	case char == 'j', key == keyboard.KeyArrowDown, key == keyboard.KeyTab:
		if self.activePos < len(self.entries)-1 {
			self.activePos++
		} else {
			self.activePos = 0
		}
	case char == 'k', key == keyboard.KeyArrowUp:
		if self.activePos > 0 {
			self.activePos--
		} else {
			self.activePos = len(self.entries) - 1
		}
	case key == keyboard.KeyEnter, char == 's':
		self.status = modalSelectStatusSelected
	case key == keyboard.KeyEsc:
		self.status = modalSelectStatusAborted
	}
	return nil
}

// print renders the dialog in its current status to the CLI.
func (self *modalSelect) print() {
	if self.status == modalSelectStatusNew {
		self.status = modalSelectStatusSelecting
	} else {
		cursor.Up(len(self.entries))
	}
	for e, entry := range self.entries {
		if e == self.initialPos && e == self.activePos { //nolint:gocritic
			self.activeColor.Println(self.initialCursor + entry.Text)
		} else if e == self.initialPos {
			self.initialColor.Println(self.initialCursor + entry.Text)
		} else if e == self.activePos {
			self.activeColor.Println(self.activeCursor + entry.Text)
		} else {
			fmt.Println(strings.Repeat(" ", len(self.activeCursor)) + entry.Text)
		}
	}
}

// selectedValue provides the value selected by the user.
func (self *modalSelect) selectedValue() string {
	return self.entries[self.activePos].Value
}
