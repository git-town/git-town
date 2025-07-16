package dialogcomponents

import (
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// SendInputs sends the given keystrokes to the given bubbletea program.
func SendInputs(input Option[TestInput], program *tea.Program) {
	if input, has := input.Get(); has {
		go func() {
			for _, msg := range input {
				program.Send(msg)
			}
		}()
	}
}
