package dialogcomponents

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcolors"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
)

func TextField(args TextFieldArgs) (string, dialogdomain.Exit, error) {
	textInput := textinput.New()
	textInput.SetValue(args.ExistingValue)
	textInput.Prompt = args.Prompt
	textInput.Focus()
	model := textFieldModel{
		colors:    dialogcolors.NewDialogColors(),
		help:      args.Help,
		status:    list.StatusActive,
		textInput: textInput,
		title:     args.Title,
	}
	program := tea.NewProgram(model)
	SendInputs(args.DialogName, args.Inputs.Next(), program)
	dialogResult, err := program.Run()
	result := dialogResult.(textFieldModel)
	return result.textInput.Value(), result.status == list.StatusExit, err
}

type TextFieldArgs struct {
	DialogName    string
	ExistingValue string
	Help          string
	Inputs        Inputs
	Prompt        string
	Title         string
}

type textFieldModel struct {
	colors    dialogcolors.DialogColors // colors to use for help text
	help      string
	status    list.Status
	textInput textinput.Model
	title     string
}

func (self textFieldModel) Init() tea.Cmd {
	return nil
}

func (self textFieldModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type { //nolint:exhaustive
		case tea.KeyEnter:
			self.status = list.StatusDone
			return self, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			self.status = list.StatusExit
			return self, tea.Quit
		}
	case error:
		panic(msg.Error())
	}
	var cmd tea.Cmd
	self.textInput, cmd = self.textInput.Update(msg)
	return self, cmd
}

func (self textFieldModel) View() string {
	if self.status != list.StatusActive {
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
