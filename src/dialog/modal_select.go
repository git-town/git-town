package dialog

import (
	"fmt"
	"strings"

	"atomicgo.dev/cursor"
	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
)

// ModalSelect allows the user to select a value from the given entries.
// Entries can be arbitrarily formatted.
// The given initial value is preselected.
func ModalSelect(entries ModalEntries, initialValue string) (*string, error) {
	initialPos := entries.IndexOfValue(initialValue)
	if initialPos == nil {
		return nil, fmt.Errorf("given initial value %q not in given entries", initialValue)
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
	activeColor   *color.Color      // color with which to print the currently selected line
	activeCursor  string            // text that gets prepended to the currently selected row
	activePos     int               // index of the currently selected row
	entries       ModalEntries      // the entries to display
	initialColor  *color.Color      // color with which to print the initially selected value
	initialCursor string            // cursor at the initial entry
	initialPos    int               // index of the initially selected value
	status        modalSelectStatus // the current status of this ModalInput instance
}

// Display shows the dialog and lets the user select an entry.
// Returns the selected value or nil if the user aborted the dialog.
func (mi *modalSelect) Display() (*string, error) {
	cursor.Hide()
	defer cursor.Show()
	err := keyboard.Open()
	if err != nil {
		return nil, err
	}
	defer keyboard.Close()
	mi.print()
	for mi.status == modalSelectStatusSelecting {
		err := mi.handleInput()
		if err != nil {
			return nil, err
		}
		mi.print()
	}
	if mi.status == modalSelectStatusAborted {
		return nil, nil //nolint:nilnil
	}
	selectedValue := mi.selectedValue()
	return &selectedValue, nil
}

// print renders the dialog in its current status to the CLI.
func (mi *modalSelect) print() {
	if mi.status == modalSelectStatusNew {
		mi.status = modalSelectStatusSelecting
	} else {
		cursor.Up(len(mi.entries))
	}
	cursorSpace := strings.Repeat(" ", len(mi.activeCursor))
	for e, entry := range mi.entries {
		if e == mi.initialPos && e == mi.activePos { //nolint:gocritic
			mi.activeColor.Println(mi.initialCursor + entry.Text)
		} else if e == mi.initialPos {
			mi.initialColor.Println(mi.initialCursor + entry.Text)
		} else if e == mi.activePos {
			mi.activeColor.Println(mi.activeCursor + entry.Text)
		} else {
			fmt.Println(cursorSpace + entry.Text)
		}
	}
}

// handleInput waits for keyboard input and updates the dialog state.
func (mi *modalSelect) handleInput() error {
	char, key, err := keyboard.GetSingleKey()
	if err != nil {
		return err
	}
	switch {
	case char == 'j', key == keyboard.KeyArrowDown, key == keyboard.KeyTab:
		if mi.activePos < len(mi.entries)-1 {
			mi.activePos++
		} else {
			mi.activePos = 0
		}
	case char == 'k', key == keyboard.KeyArrowUp:
		if mi.activePos > 0 {
			mi.activePos--
		} else {
			mi.activePos = len(mi.entries) - 1
		}
	case key == keyboard.KeyEnter, char == 's':
		mi.status = modalSelectStatusSelected
	case key == keyboard.KeyEsc:
		mi.status = modalSelectStatusAborted
	}
	return nil
}

// selectedValue provides the value selected by the user.
func (mi *modalSelect) selectedValue() string {
	return mi.entries[mi.activePos].Value
}

// ModalEntry contains one of the many entries that the user can choose from.
type ModalEntry struct {
	Text  string // the text to display
	Value string // the return value
}

// ModalEntries is a collection of ModalEntry.
type ModalEntries []ModalEntry

// IndexOfValue provides the index of the entry with the given value,
// or nil if the given value is not in the list.
func (mes ModalEntries) IndexOfValue(value string) *int {
	for e, entry := range mes {
		if entry.Value == value {
			return &e
		}
	}
	return nil
}

// modalSelectStatus represents the different states that a modalSelect instance can be in.
type modalSelectStatus int

const (
	modalSelectStatusNew modalSelectStatus = iota
	modalSelectStatusSelecting
	modalSelectStatusSelected
	modalSelectStatusAborted
)
