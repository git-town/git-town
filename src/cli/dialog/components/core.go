// Package dialog contains shared components to build full dialog screens out of.
package components

import tea "github.com/charmbracelet/bubbletea"

func SendInputs(inputs TestInput, program *tea.Program) {
	if len(inputs) > 0 {
		go func() {
			for _, input := range inputs {
				program.Send(input)
			}
		}()
	}
}
