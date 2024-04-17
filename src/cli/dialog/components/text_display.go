package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v14/src/cli/colors"
	"github.com/git-town/git-town/v14/src/cli/dialog/components/list"
)

func TextDisplay(title, text string, inputs TestInput) (bool, error) {
	model := textDisplayModel{
		colors: colors.CreateColors(),
		status: list.StatusActive,
		text:   text,
		title:  title,
	}
	program := tea.NewProgram(model)
	SendInputs(inputs, program)
	dialogResult, err := program.Run()
	if err != nil {
		return false, err
	}
	result := dialogResult.(textDisplayModel) //nolint:forcetypeassert
	return result.status == list.StatusAborted, nil
}

type textDisplayModel struct {
	colors colors.DialogColors
	status list.Status
	text   string
	title  string
}

func (self textDisplayModel) Init() tea.Cmd {
	return nil
}

func (self textDisplayModel) Update(msg tea.Msg) (model tea.Model, cmd tea.Cmd) { //nolint:ireturn
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type { //nolint:exhaustive
		case tea.KeyEnter:
			self.status = list.StatusDone
			return self, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			self.status = list.StatusAborted
			return self, tea.Quit
		}
		switch msg.String() {
		case "o":
			self.status = list.StatusDone
			return self, tea.Quit
		case "q":
			self.status = list.StatusAborted
			return self, tea.Quit
		}
	case error:
		panic(msg.Error())
	}
	return self, cmd
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
