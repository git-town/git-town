package dialog

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func TextInput(existingValue string, help string, prompt string, testInput TestInput) (string, bool, error) {
	textInput := textinput.New()
	textInput.SetValue(existingValue)
	textInput.Prompt = prompt
	textInput.Focus()
	model := textInputModel{
		textInput: textInput,
		colors:    createColors(),
		help:      help,
		status:    StatusActive,
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
	return result.textInput.Value(), result.status == StatusAborted, nil
}

type textInputModel struct {
	textInput textinput.Model
	colors    dialogColors // colors to use for help text
	help      string
	status    status
}

func (m textInputModel) Init() tea.Cmd {
	return nil
}

func (self textInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			self.status = StatusDone
			return self, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			self.status = StatusAborted
			return self, tea.Quit
		}
	case error:
		panic(msg.Error())
	}
	var cmd tea.Cmd
	self.textInput, cmd = self.textInput.Update(msg)
	return self, cmd
}

func (self textInputModel) View() string {
	if self.status != StatusActive {
		return ""
	}
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
