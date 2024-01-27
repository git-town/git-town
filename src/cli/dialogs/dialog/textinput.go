package dialog

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func TextInput(existingValue string, help string, testInput TestInput) (string, bool, error) {
	textInput := textinput.New()
	textInput.SetValue(existingValue)
	textInput.Prompt = "Your GitHub token: "
	textInput.Focus()
	model := textInputModel{
		textInput: textInput,
		colors:    createColors(),
		help:      help,
		aborted:   false,
	}
	program := tea.NewProgram(model)
	if len(testInput) > 0 {
		go func() {
			for _, input := range testInput {
				program.Send(input)
			}
		}()
	}
	dialogResult, err := program.Run()
	if err != nil {
		return existingValue, false, err
	}
	result := dialogResult.(textInputModel) //nolint:forcetypeassert
	return result.textInput.Value(), result.aborted, nil
}

type textInputModel struct {
	textInput textinput.Model
	colors    dialogColors // colors to use for help text
	help      string
	aborted   bool
}

func (m textInputModel) Init() tea.Cmd {
	return nil
}

func (m textInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.aborted = true
			return m, tea.Quit
		}
	case error:
		panic(msg.Error())
	}
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (self textInputModel) View() string {
	result := strings.Builder{}
	result.WriteString(self.help)
	result.WriteString(self.textInput.View())
	result.WriteString("\n\n  ")
	// accept
	result.WriteString(self.colors.HelpKey.Styled("enter"))
	result.WriteString(self.colors.Help.Styled(" accept   "))
	// abort
	result.WriteString(self.colors.HelpKey.Styled("esc"))
	result.WriteString(self.colors.Help.Styled("/"))
	result.WriteString(self.colors.HelpKey.Styled("ctrl-c"))
	result.WriteString(self.colors.Help.Styled(" abort"))
	return result.String()
}
