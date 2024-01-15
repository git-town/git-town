package dialog

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

func SwitchBranch(localBranches gitdomain.LocalBranchNames, initialBranch gitdomain.LocalBranchName) (gitdomain.LocalBranchName, error) {
	dialogData := SwitchModel{
		BubbleList:    newBubbleList(localBranches.Strings(), initialBranch.String()),
		InitialBranch: initialBranch.String(),
	}
	dialogProcess := tea.NewProgram(dialogData, tea.WithOutput(os.Stderr))
	dialogResult, err := dialogProcess.Run()
	if err != nil {
		return "", err
	}
	result := dialogResult.(SwitchModel) //nolint:forcetypeassert
	selectedBranch := gitdomain.NewLocalBranchName(result.BubbleList.selectedEntry())
	return selectedBranch, nil
}

type SwitchModel struct {
	BubbleList
	InitialBranch string // name of the currently checked out branch
}

func (self SwitchModel) Init() tea.Cmd {
	return nil
}

func (self SwitchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:ireturn
	keyMsg, isKeyMsg := msg.(tea.KeyMsg)
	if !isKeyMsg {
		return self, nil
	}
	if handled, code := self.BubbleList.handleKey(keyMsg); handled {
		return self, code
	}
	if keyMsg.Type == tea.KeyEnter {
		return self, tea.Quit
	}
	if keyMsg.String() == "o" {
		return self, tea.Quit
	}
	return self, nil
}

func (self SwitchModel) View() string {
	s := strings.Builder{}
	for i, branch := range self.Entries {
		switch {
		case i == self.Cursor:
			s.WriteString(self.Colors.selection.Styled("> " + branch))
		case branch == self.InitialBranch:
			s.WriteString(self.Colors.initial.Styled("* " + branch))
		default:
			s.WriteString("  " + branch)
		}
		s.WriteRune('\n')
	}
	s.WriteString("\n\n  ")
	// up
	s.WriteString(self.Colors.helpKey.Styled("↑"))
	s.WriteString(self.Colors.help.Styled("/"))
	s.WriteString(self.Colors.helpKey.Styled("k"))
	s.WriteString(self.Colors.help.Styled(" up   "))
	// down
	s.WriteString(self.Colors.helpKey.Styled("↓"))
	s.WriteString(self.Colors.help.Styled("/"))
	s.WriteString(self.Colors.helpKey.Styled("j"))
	s.WriteString(self.Colors.help.Styled(" down   "))
	// accept
	s.WriteString(self.Colors.helpKey.Styled("enter"))
	s.WriteString(self.Colors.help.Styled("/"))
	s.WriteString(self.Colors.helpKey.Styled("o"))
	s.WriteString(self.Colors.help.Styled(" accept   "))
	// abort
	s.WriteString(self.Colors.helpKey.Styled("esc"))
	s.WriteString(self.Colors.help.Styled("/"))
	s.WriteString(self.Colors.helpKey.Styled("q"))
	s.WriteString(self.Colors.help.Styled(" abort"))
	return s.String()
}
