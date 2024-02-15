package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func TextField(args TextFieldArgs) (string, bool, error) {
	textInput := textinput.New()
	textInput.SetValue(args.ExistingValue)
	textInput.Prompt = args.Prompt
	textInput.Focus()
	model := textFieldModel{
		colors:    createColors(),
		help:      args.Help,
		status:    StatusActive,
		textInput: textInput,
		title:     args.Title,
	}
	program := tea.NewProgram(model)
	SendInputs(args.TestInput, program)
	dialogResult, err := program.Run()
	if err != nil {
		return args.ExistingValue, false, err
	}
	result := dialogResult.(textFieldModel) //nolint:forcetypeassert
	return result.textInput.Value(), result.status == StatusAborted, nil
}

type TextFieldArgs struct {
	ExistingValue string
	Help          string
	Prompt        string
	TestInput     TestInput
	Title         string
}

type textFieldModel struct {
	colors    dialogColors // colors to use for help text
	help      string
	status    status
	textInput textinput.Model
	title     string
}

func (self textFieldModel) Init() tea.Cmd {
	return nil
}

func (self textFieldModel) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) { //nolint:ireturn
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type { //nolint:exhaustive
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
	self.textInput, cmd = self.textInput.Update(msg)
	return self, cmd
}

func (self textFieldModel) View() string {
	if self.status != StatusActive {
		return ""
	}
	result := strings.Builder{}
	result.WriteRune('\n')
	result.WriteString(self.colors.Title.Styled(self.title))
	result.WriteRune('\n')
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
