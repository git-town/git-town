package dialogcomponents

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcolors"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogcomponents/list"
	"github.com/git-town/git-town/v22/internal/cli/dialog/dialogdomain"
)

func TextDisplay(title, text string, inputs Inputs, dialogName string) (dialogdomain.Exit, error) {
	model := textDisplayModel{
		colors: dialogcolors.NewDialogColors(),
		status: list.StatusActive,
		text:   text,
		title:  title,
	}
	program := tea.NewProgram(model)
	SendInputs(dialogName, inputs.Next(), program)
	dialogResult, err := program.Run()
	result := dialogResult.(textDisplayModel)
	return result.status == list.StatusExit, err
}

type textDisplayModel struct {
	colors dialogcolors.DialogColors
	status list.Status
	text   string
	title  string
}

func (self textDisplayModel) Init() tea.Cmd {
	return nil
}

func (self textDisplayModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
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
		switch msg.String() {
		case "o":
			self.status = list.StatusDone
			return self, tea.Quit
		case "q":
			self.status = list.StatusExit
			return self, tea.Quit
		}
	case error:
		panic(msg.Error())
	}
	return self, nil
}

func (self textDisplayModel) View() string {
	if self.status != list.StatusActive {
		return ""
	}
	result := strings.Builder{}
	result.WriteRune('\n')
	result.WriteString(self.colors.Title.Styled(self.title))
	result.WriteRune('\n')
	result.WriteString(self.text)
	result.WriteString("\n\n  ")
	// accept
	result.WriteString(self.colors.HelpKey.Styled("o"))
	result.WriteString(self.colors.Help.Styled("/"))
	result.WriteString(self.colors.HelpKey.Styled("enter"))
	result.WriteString(self.colors.Help.Styled(" continue   "))
	// abort
	result.WriteString(self.colors.HelpKey.Styled("q"))
	result.WriteString(self.colors.Help.Styled("/"))
	result.WriteString(self.colors.HelpKey.Styled("esc"))
	result.WriteString(self.colors.Help.Styled("/"))
	result.WriteString(self.colors.HelpKey.Styled("ctrl-c"))
	result.WriteString(self.colors.Help.Styled(" abort"))
	return result.String()
}
