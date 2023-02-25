package dialog

import (
	"fmt"
	"strings"

	"github.com/fatih/color"

	"atomicgo.dev/cursor"
	"github.com/eiannone/keyboard"
)

// ModalInput allows selecting a value using VIM keybindings.
type ModalInput struct {
	// the entries to display
	entries []ModalEntry

	// cursorPos contains the index of the currently selected row.
	cursorPos int

	// cursorText contains the text of the cursor, including color codes.
	cursorText string

	activeLineColor *color.Color

	// Result contains the result that the user has selected,
	// or nil if no selection has taken place yet.
	status modalInputStatus
}

func NewModalInput(entries []ModalEntry, cursorText string, initialValue string) (*ModalInput, func(), error) {
	cursor.Hide()
	if err := keyboard.Open(); err != nil {
		return nil, nil, err
	}
	cursorPos := 0
	for e, entry := range entries {
		if entry.Value == initialValue {
			cursorPos = e
			break
		}
	}
	input := ModalInput{
		activeLineColor: color.New(color.FgCyan, color.Bold),
		entries:         entries,
		cursorPos:       cursorPos,
		cursorText:      cursorText,
		status:          modalInputStatusNew,
	}
	return &input, input.cleanup, nil
}

func (mi *ModalInput) Display() (*string, error) {
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

func (mi *ModalInput) cleanup() {
	cursor.Show()
	keyboard.Close()
}

// Display displays this dialog.
func (mi *ModalInput) print() {
	if mi.status == modalInputStatusNew {
		mi.status = modalInputStatusSelecting
	} else {
		cursor.Up(len(mi.entries))
	}
	cursorSpace := strings.Repeat(" ", len(mi.cursorText))
	for e := range mi.entries {
		if e == mi.cursorPos {
			mi.activeLineColor.Println(mi.cursorText + mi.entries[e].Text)
		} else {
			fmt.Println(cursorSpace + mi.entries[e].Text)
		}
	}
}

// Process waits for keyboard input, updates the dialog state, and re-draws the dialog.
func (mi *ModalInput) handleInput() error {
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

func (mi *ModalInput) selectedValue() string {
	return mi.entries[mi.cursorPos].Value
}

type ModalEntry struct {
	// the text to display
	Text string

	// the return value
	Value string
}

type modalInputStatus int

const (
	modalInputStatusNew modalInputStatus = iota
	modalInputStatusSelecting
	modalInputStatusSelected
	modalInputStatusAborted
)
