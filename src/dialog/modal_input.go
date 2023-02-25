package dialog

import (
	"fmt"
	"strings"

	"github.com/eiannone/keyboard"
)

// ModalInput allows selecting a value using VIM keybindings.
type ModalInput struct {
	// the entries to display
	Entries []ModalEntry

	// CursorPos contains the index of the currently selected row.
	CursorPos uint8

	// CursorText contains the text of the cursor, including color codes.
	CursorText string

	// Result contains the result that the user has selected,
	// or nil if no selection has taken place yet.
	Status ModalInputStatus
}

// Display displays this dialog.
func (mi *ModalInput) Display() {
	cursorSpace := strings.Repeat(" ", len(mi.CursorText))
	for e := range mi.Entries {
		if e == int(mi.CursorPos) {
			fmt.Println(mi.CursorText + mi.Entries[e].Text)
		} else {
			fmt.Println(cursorSpace + mi.Entries[e].Text)
		}
	}
}

func (mi *ModalInput) HandleInput() error {
	char, key, err := keyboard.GetSingleKey()
	if err != nil {
		return err
	}
	if char == 'j' || key == keyboard.KeyArrowDown {
		mi.CursorPos += 1
	} else if char == 'k' || key == keyboard.KeyArrowUp {
		mi.CursorPos -= 1
	} else if key == keyboard.KeyEnter || key == keyboard.KeySpace {
		mi.Status = ModalInputStatusSelected
	} else if key == keyboard.KeyEsc {
		mi.Status = ModalInputStatusAborted
	}
	return nil
}

func (mi *ModalInput) SelectedValue() string {
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
	ModalInputStatusSelecting ModalInputStatus = iota
	ModalInputStatusSelected
	ModalInputStatusAborted
)
