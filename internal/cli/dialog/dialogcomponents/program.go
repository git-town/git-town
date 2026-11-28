package dialogcomponents

import (
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/git-town/git-town/v24/pkg/prelude"
)

// NewProgram creates a Bubble Tea program. Scripted dialog input disables the renderer
// so that transient dialog frames don't pollute command output captured by end-to-end tests.
func NewProgram(model tea.Model, input Option[Input]) *tea.Program {
	if input.IsSome() {
		return tea.NewProgram(model, tea.WithoutRenderer())
	}
	return tea.NewProgram(model)
}
