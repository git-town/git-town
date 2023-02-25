package dialog

import (
	"fmt"
	"strings"

	"atomicgo.dev/cursor"
	"github.com/eiannone/keyboard"
)

// ModalInput allows selecting a value using VIM keybindings.
type ModalInput struct {
	// the entries to display
	Entries []ModalEntry

	// CursorPos contains the index of the currently selected row.
	CursorPos int

	// CursorText contains the text of the cursor, including color codes.
	CursorText string

	// Result contains the result that the user has selected,
	// or nil if no selection has taken place yet.
	Status ModalInputStatus
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
		Entries:    entries,
		CursorPos:  cursorPos,
		CursorText: cursorText,
		Status:     ModalInputStatusNew,
	}
	return &input, input.Cleanup, nil
}

func (mi *ModalInput) Cleanup() {
	cursor.Show()
	keyboard.Close()
}
func (mi *ModalInput) Display() (*string, error) {
	mi.print()
	for mi.Status == ModalInputStatusSelecting {
		err := mi.handleInput()
		if err != nil {
			return nil, err
		}
		mi.print()
	}
	if mi.Status == ModalInputStatusAborted {
		return nil, nil
	}
	selectedValue := mi.selectedValue()
	return &selectedValue, nil
}

// Display displays this dialog.
func (mi *ModalInput) print() {
	if mi.Status == ModalInputStatusNew {
		mi.Status = ModalInputStatusSelecting
	} else {
		cursor.Up(len(mi.Entries))
	}
	cursorSpace := strings.Repeat(" ", len(mi.CursorText))
	for e := range mi.Entries {
		if e == int(mi.CursorPos) {
			fmt.Println(mi.CursorText + mi.Entries[e].Text)
		} else {
			fmt.Println(cursorSpace + mi.Entries[e].Text)
		}
	}
}

// Process waits for keyboard input, updates the dialog state, and re-draws the dialog.
func (mi *ModalInput) handleInput() error {
	char, key, err := keyboard.GetSingleKey()
	if err != nil {
		return err
	}
	if char == 'j' || key == keyboard.KeyArrowDown || key == keyboard.KeyTab {
		if mi.CursorPos < len(mi.Entries)-1 {
			mi.CursorPos += 1
		} else {
			mi.CursorPos = 0
		}
	} else if char == 'k' || key == keyboard.KeyArrowUp {
		if mi.CursorPos > 0 {
			mi.CursorPos -= 1
		} else {
			mi.CursorPos = len(mi.Entries) - 1
		}
	} else if key == keyboard.KeyEnter || char == 's' {
		mi.Status = ModalInputStatusSelected
	} else if key == keyboard.KeyEsc {
		mi.Status = ModalInputStatusAborted
	}
	return nil
}

func (mi *ModalInput) selectedValue() string {
	return mi.Entries[mi.CursorPos].Value
}

type ModalEntry struct {
	// the text to display
	Text string

	// the return value
	Value string
}

type ModalInputStatus int

const (
	ModalInputStatusNew ModalInputStatus = iota
	ModalInputStatusSelecting
	ModalInputStatusSelected
	ModalInputStatusAborted
)
