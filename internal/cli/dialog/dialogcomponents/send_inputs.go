package dialogcomponents

import tea "github.com/charmbracelet/bubbletea"

// SendInputs sends the given keystrokes to the given bubbletea program.
func SendInputs(inputs TestInput, program *tea.Program) {
	if len(inputs) > 0 {
		go func() {
			for _, input := range inputs {
				program.Send(input)
			}
		}()
	}
}
