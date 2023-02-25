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
func ModalSelect(entries ModalEntries, cursorText string, initialValue string) (*string, error) {
	cursorPos := entries.IndexOfValue(initialValue)
	if cursorPos == nil {
		return nil, fmt.Errorf("given initial value %q not in given entries", initialValue)
	}
	input := modalSelect{
		activeColor: color.New(color.FgCyan, color.Bold),
		entries:     entries,
		cursorPos:   *cursorPos,
		cursorText:  cursorText,
		status:      modalInputStatusNew,
	}
	return input.Display()
}

// modalSelect allows selecting a value from a list using VIM keybindings.
type modalSelect struct {
	activeColor *color.Color     // color with which to print the currently selected line
	cursorPos   int              // index of the currently selected row
	cursorText  string           // text that gets prepended to the currently selected row
	entries     ModalEntries     // the entries to display
	status      modalInputStatus // the current status of this ModalInput instance
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
	for mi.status == modalInputStatusSelecting {
		err := mi.handleInput()
		if err != nil {
			return nil, err
		}
		mi.print()
	}
	if mi.status == modalInputStatusAborted {
		return nil, nil //nolint:nilnil
	}
	selectedValue := mi.selectedValue()
	return &selectedValue, nil
}

// print renders the dialog in its current status to the CLI.
func (mi *modalSelect) print() {
	if mi.status == modalInputStatusNew {
		mi.status = modalInputStatusSelecting
	} else {
		cursor.Up(len(mi.entries))
	}
	cursorSpace := strings.Repeat(" ", len(mi.cursorText))
	for e, entry := range mi.entries {
		if e == mi.cursorPos {
			mi.activeColor.Println(mi.cursorText + entry.Text)
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
		if mi.cursorPos < len(mi.entries)-1 {
			mi.cursorPos++
		} else {
			mi.cursorPos = 0
		}
	case char == 'k', key == keyboard.KeyArrowUp:
		if mi.cursorPos > 0 {
			mi.cursorPos--
		} else {
			mi.cursorPos = len(mi.entries) - 1
		}
	case key == keyboard.KeyEnter, char == 's':
		mi.status = modalInputStatusSelected
	case key == keyboard.KeyEsc:
		mi.status = modalInputStatusAborted
	}
	return nil
}

// selectedValue provides the value selected by the user.
func (mi *modalSelect) selectedValue() string {
	return mi.entries[mi.cursorPos].Value
}

type ModalEntry struct {
	// the text to display
	Text string

	// the return value
	Value string
}

type ModalEntries []ModalEntry

// IndexOfValue provides the index of the entry with the given value.
func (mes ModalEntries) IndexOfValue(value string) *int {
	for e, entry := range mes {
		if entry.Value == value {
			return &e
		}
	}
	return nil
}

type modalInputStatus int

const (
	modalInputStatusNew modalInputStatus = iota
	modalInputStatusSelecting
	modalInputStatusSelected
	modalInputStatusAborted
)
