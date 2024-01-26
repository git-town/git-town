package enter

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/cli/dialogs/dialog"
)

func textInput(existingValue string, help string, placeholder string, testInput dialog.TestInput) (string, bool, error) {
	ti := textinput.New()
	ti.SetValue(existingValue)
	ti.Placeholder = placeholder
	ti.Focus()
	model := textInputModel{
		textInput: ti,
	}
	program := tea.NewProgram(model)
	dialogResult, err := program.Run()
	if err != nil {
		return existingValue, false, err
	}
	result := dialogResult.(textInputModel) //nolint:forcetypeassert
	return result.textInput.Value(), result.aborted, nil
}

type textInputModel struct {
	textInput textinput.Model
	help      string
	aborted   bool
	err       error
}

func (m textInputModel) Init() tea.Cmd {
	return nil
}

func (m textInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			m.aborted = true
			return m, tea.Quit
		}
	case error:
		m.err = msg
	}
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m textInputModel) View() string {
	result := strings.Builder{}
	result.WriteString(m.help)
	result.WriteString(m.textInput.View())
	result.WriteString("\n\n(esc to quit)\n")
	return result.String()
}
